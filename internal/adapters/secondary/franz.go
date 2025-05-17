package secondary

import (
	"context"
	"log/slog"

	"github.com/nxdir-s/IdleEngine/internal/adapters/secondary/franz"
	"github.com/nxdir-s/IdleEngine/internal/util"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	MaxPollFetches int = 1000
)

type ErrProtoMarshal struct {
	err error
}

func (e *ErrProtoMarshal) Error() string {
	return "failed to marshal protobuf: " + e.err.Error()
}

type ErrCreateTopic struct {
	err error
}

func (e *ErrCreateTopic) Error() string {
	return "failed to create kafka topic: " + e.err.Error()
}

type ErrListTopics struct {
	err error
}

func (e *ErrListTopics) Error() string {
	return "failed to list kafka topics: " + e.err.Error()
}

type FranzAdapterOpt func(a *FranzAdapter) error

func WithConsumer(topic string, groupname string, brokers []string) FranzAdapterOpt {
	return func(a *FranzAdapter) error {
		client, err := kgo.NewClient(
			kgo.SeedBrokers(brokers...),
			kgo.ConsumerGroup(groupname),
			kgo.ConsumeTopics(topic),
			kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
			kgo.DisableAutoCommit(),
			kgo.BlockRebalanceOnPoll(),
		)
		if err != nil {
			return err
		}

		a.topic = topic
		a.client = client
		a.groupName = groupname

		return nil
	}
}

func WithProducer(topic string, brokers []string) FranzAdapterOpt {
	return func(a *FranzAdapter) error {
		client, err := kgo.NewClient(
			kgo.SeedBrokers(brokers...),
		)
		if err != nil {
			return err
		}

		a.topic = topic
		a.client = client

		return nil
	}
}

type FranzAdapter struct {
	client    *kgo.Client
	logger    *slog.Logger
	tracer    trace.Tracer
	topic     string
	groupName string
}

// NewFranzAdapter creates a new kafka adapter. Adapters should be configured to either produce or consume, but not both
func NewFranzAdapter(logger *slog.Logger, tracer trace.Tracer, opts ...FranzAdapterOpt) (*FranzAdapter, error) {
	adapter := &FranzAdapter{
		logger: logger,
		tracer: tracer,
	}

	for _, opt := range opts {
		if err := opt(adapter); err != nil {
			return nil, err
		}
	}

	return adapter, nil
}

// Send sends the supplied kafka record to the configured topic
func (a *FranzAdapter) Send(ctx context.Context, record protoreflect.ProtoMessage) error {
	if a.client == nil {
		a.logger.Error("nil client in FranzAdapter")
		return nil
	}

	ctx, span := a.tracer.Start(ctx, "send "+a.topic,
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination.name", a.topic),
			attribute.String("messaging.operation.name", "send"),
			attribute.String("messaging.operation.type", "send"),
		),
	)
	defer span.End()

	data, err := proto.Marshal(record)
	if err != nil {
		a.logger.Error("error encoding record",
			slog.String("err", err.Error()),
			slog.String("topic", a.topic),
		)

		err = &ErrProtoMarshal{err}
		util.RecordError(span, "error encoding "+a.topic+" record", err)

		return err
	}

	a.logger.Debug("sending kafka record")

	a.client.Produce(ctx, &kgo.Record{Topic: a.topic, Value: data}, func(_ *kgo.Record, err error) {
		if err != nil {
			a.logger.Error("record had a produce error", slog.Any("err", err))
		}
	})

	return nil
}

// Consume uses the supplied consumer to process kafka records from the configured topic
func (a *FranzAdapter) Consume(ctx context.Context, consumer franz.Consumer) {
	if a.client == nil {
		a.logger.Error("nil client in FranzAdapter")
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			ctx, span := a.tracer.Start(ctx, "receive "+a.topic,
				trace.WithSpanKind(trace.SpanKindClient),
				trace.WithAttributes(
					attribute.String("messaging.system", "kafka"),
					attribute.String("messaging.consumer.group.name", a.groupName),
					attribute.String("messaging.operation.name", "receive"),
					attribute.String("messaging.operation.type", "receive"),
				),
			)

			fetches := a.client.PollRecords(ctx, MaxPollFetches)
			span.End()

			ctx, span = a.tracer.Start(ctx, "process "+a.topic,
				trace.WithSpanKind(trace.SpanKindConsumer),
				trace.WithAttributes(
					attribute.String("messaging.system", "kafka"),
					attribute.String("messaging.consumer.group.name", a.groupName),
					attribute.String("messaging.operation.name", "process"),
					attribute.String("messaging.operation.type", "process"),
				),
			)

			if errors := fetches.Errors(); len(errors) > 0 {
				for _, e := range errors {
					if e.Err == context.Canceled {
						a.logger.Error("received interrupt", slog.String("err", e.Err.Error()))

						util.RecordError(span, "received interrupt", e.Err)
						span.End()

						return
					}

					a.logger.Error("poll error", slog.String("err", e.Err.Error()))
				}
			}

			iter := fetches.RecordIter()
			for !iter.Done() {
				record := iter.Next()

				if err := consumer.Process(ctx, record); err != nil {
					a.logger.Error("error processing "+a.topic+" record", slog.String("err", err.Error()))

					util.RecordError(span, "error processing "+a.topic+" record", err)
					span.End()

					return // return or continue?
				}

				a.logger.Info("consumed record", slog.String("topic", a.topic))
			}

			if err := a.client.CommitUncommittedOffsets(ctx); err != nil {
				if err == context.Canceled {
					a.logger.Error("received interrupt", slog.String("err", err.Error()))

					util.RecordError(span, "received interrupt", err)
					span.End()

					return
				}

				a.logger.Error("unable to commit offsets", slog.String("err", err.Error()))
			}

			a.client.AllowRebalance()
			span.End()
		}
	}
}

// Close closes the kafka client
func (a *FranzAdapter) Close() error {
	if a.client == nil {
		return nil
	}

	a.client.Close()

	return nil
}

// CreateTopic creates a kafka topic
func (a *FranzAdapter) CreateTopic(ctx context.Context, topic string) error {
	if a.client == nil {
		a.logger.Error("nil client in FranzAdapter")
		return nil
	}

	ctx, span := a.tracer.Start(ctx, "create.topic "+a.topic,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination.name", topic),
			attribute.String("messaging.operation.name", "create topic"),
		),
	)
	defer span.End()

	adminClient := kadm.NewClient(a.client)

	topicDetails, err := adminClient.ListTopics(ctx)
	if err != nil {
		return &ErrListTopics{err}
	}

	if topicDetails.Has(topic) {
		a.logger.Info("kafka topic already exists", slog.String("topic", topic))
		return nil
	}

	a.logger.Info("creating kafka topic", slog.String("topic", topic))

	if _, err := kadm.NewClient(a.client).CreateTopic(ctx, 1, -1, nil, topic); err != nil {
		a.logger.Error("failed to create kafka topic",
			slog.String("err", err.Error()),
			slog.String("topic", topic),
		)

		return &ErrCreateTopic{err}
	}

	return nil
}
