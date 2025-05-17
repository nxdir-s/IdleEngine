package engine

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/nxdir-s/IdleEngine/internal/core/valobj"
	"github.com/nxdir-s/IdleEngine/internal/ports"
	"github.com/nxdir-s/IdleEngine/internal/server"
	"github.com/nxdir-s/IdleEngine/protobuf"
	"github.com/nxdir-s/pipelines"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	SimulateMaxFan int = 3
	KafkaMaxFan    int = 3

	TickerInterval time.Duration = time.Second * 5
)

type GameEngine struct {
	kafka    ports.Kafka
	pool     *server.Pool
	ticker   *time.Ticker
	tracer   trace.Tracer
	logger   *slog.Logger
	sigusr1  chan os.Signal
	isPaused bool
}

// NewGameEngine creates a GameEngine
func NewGameEngine(pool *server.Pool, kafka ports.Kafka, logger *slog.Logger, tracer trace.Tracer) *GameEngine {
	return &GameEngine{
		kafka:   kafka,
		pool:    pool,
		ticker:  time.NewTicker(TickerInterval),
		sigusr1: make(chan os.Signal, 1),
		tracer:  tracer,
		logger:  logger,
	}
}

// Start runs the game engine
func (ngin *GameEngine) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ngin.sigusr1:
			ngin.isPaused = !ngin.isPaused

			switch ngin.isPaused {
			case true:
				ngin.logger.Info("engine is paused")
			case false:
				ngin.logger.Info("engine is resumed")
			}
		case t := <-ngin.ticker.C:
			if ngin.isPaused {
				break
			}

			ctx, span := ngin.tracer.Start(ctx, "engine tick")

			ngin.logger.Info("server tick")

			reply := make(chan *server.Snapshot)
			ngin.pool.Snapshot <- reply

			snapshot := <-reply

			ngin.process(ctx, snapshot.Connections)
			snapshot.Processed <- struct{}{}

			event := ngin.buildEvent(ctx, t)

			ngin.pool.Broadcast <- event
			<-event.Consumed

			span.End()
		}
	}
}

func (ngin *GameEngine) process(ctx context.Context, users map[int32]*server.Client) {
	if len(users) == 0 {
		ngin.logger.Info("0 users connected")
		return
	}

	ctx, span := ngin.tracer.Start(ctx, "engine process")
	defer span.End()

	ngin.logger.Info("processing user events", slog.Int("connections", len(users)))

	stream := pipelines.StreamMap[int32, *server.Client](ctx, users)

	fanOutChannels := pipelines.FanOut(ctx, stream, ngin.Simulate, SimulateMaxFan)
	playerEvents := pipelines.FanIn(ctx, fanOutChannels...)

	kafkaFanOut := pipelines.FanOut(ctx, playerEvents, ngin.kafka.Send, KafkaMaxFan)
	errChan := pipelines.FanIn(ctx, kafkaFanOut...)

	for err := range errChan {
		select {
		case <-ctx.Done():
			return
		default:
			if err != nil {
				ngin.logger.Error("failed to send user event", slog.String("err", err.Error()))
			}
		}
	}
}

// Simulate simulates user events
func (ngin *GameEngine) Simulate(ctx context.Context, client *server.Client) protoreflect.ProtoMessage {
	return &protobuf.UserEvent{
		Id: client.User.Id,
	}
}

func (ngin *GameEngine) buildEvent(ctx context.Context, t time.Time) *valobj.Event {
	return &valobj.Event{
		Ctx:      ctx,
		Consumed: make(chan struct{}),
		Msg: &valobj.Message{
			Value: "server tick: " + t.UTC().String(),
		},
	}
}
