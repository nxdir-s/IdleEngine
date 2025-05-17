package consumers

import (
	"context"
	"log/slog"

	"github.com/nxdir-s/IdleEngine/internal/ports"
	"github.com/nxdir-s/IdleEngine/protobuf"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/proto"
)

type UserEventsOpts func(c *UserEvents)

type UserEvents struct {
	kafka   ports.Kafka
	adapter ports.Consumer
	logger  *slog.Logger
}

func NewUserEvents(kafka ports.Kafka, adapter ports.Consumer, logger *slog.Logger, opts ...UserEventsOpts) *UserEvents {
	consumer := &UserEvents{
		kafka:   kafka,
		adapter: adapter,
		logger:  logger,
	}

	for _, opt := range opts {
		opt(consumer)
	}

	return consumer
}

func (c *UserEvents) Start(ctx context.Context) {
	c.logger.Info("starting userevents consumer")
	c.kafka.Consume(ctx, c)
}

func (c *UserEvents) Process(ctx context.Context, record *kgo.Record) error {
	var event protobuf.UserEvent
	if err := proto.Unmarshal(record.Value, &event); err != nil {
		return err
	}

	return c.adapter.ProcessUserEvent(ctx, &event)
}

func (c *UserEvents) Close() {
	c.logger.Info("closing userevents consumer")

	if err := c.kafka.Close(); err != nil {
		c.logger.Error("error closing userevents consumer",
			slog.String("err", err.Error()),
		)
	}
}
