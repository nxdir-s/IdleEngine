package primary

import (
	"context"
	"log/slog"

	"github.com/nxdir-s/IdleEngine/internal/ports"
	"github.com/nxdir-s/IdleEngine/protobuf"
	"go.opentelemetry.io/otel/trace"
)

type ErrUserEvent struct {
	err error
}

func (e *ErrUserEvent) Error() string {
	return "failed to process user event: " + e.err.Error()
}

type ConsumerAdapter struct {
	events ports.Events
	tracer trace.Tracer
	logger *slog.Logger
}

func NewConsumerAdapter(events ports.Events, logger *slog.Logger, tracer trace.Tracer) *ConsumerAdapter {
	return &ConsumerAdapter{
		events: events,
		tracer: tracer,
		logger: logger,
	}
}

func (a *ConsumerAdapter) ProcessUserEvent(ctx context.Context, event *protobuf.UserEvent) error {
	ctx, span := a.tracer.Start(ctx, "consumer userevent")
	defer span.End()

	a.logger.Info("consumed user event", slog.Int("user_id", int(event.Id)))

	// if err := a.events.HandleUserEvent(ctx, event); err != nil {
	// 	return &ErrUserEvent{err}
	// }

	return nil
}
