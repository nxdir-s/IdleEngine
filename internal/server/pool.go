package server

import (
	"context"
	"log/slog"

	"github.com/nxdir-s/IdleEngine/internal/core/valobj"
	"github.com/nxdir-s/pipelines"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sys/unix"
)

const (
	MaxSendFan int = 3
)

type EpollEvent struct {
	Events []unix.EpollEvent
	Resp   chan []*Client
}

type Connections map[int32]*Client

type Snapshot struct {
	Connections Connections
	Processed   chan struct{}
}

type Pool struct {
	Connections Connections

	Register    chan *Client
	Remove      chan int32
	Broadcast   chan *valobj.Event
	Snapshot    chan chan *Snapshot
	EpollEvents chan *EpollEvent

	counter int32
	tracer  trace.Tracer
	logger  *slog.Logger
}

// NewPool creates a new client pool
func NewPool(ctx context.Context, logger *slog.Logger, tracer trace.Tracer) *Pool {
	return &Pool{
		Connections: make(map[int32]*Client),
		Register:    make(chan *Client),
		Remove:      make(chan int32),
		Broadcast:   make(chan *valobj.Event),
		Snapshot:    make(chan chan *Snapshot),
		EpollEvents: make(chan *EpollEvent),
		tracer:      tracer,
		logger:      logger,
	}
}

// Start runs the client pool
func (p *Pool) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case client := <-p.Register:
			_, span := p.tracer.Start(ctx, "pool register")

			// using counter for testing
			p.counter++
			client.User.Id = p.counter

			p.Connections[client.Fd] = client

			p.logger.Info("added client to pool", slog.Int("connections", len(p.Connections)))
			span.End()
		case fd := <-p.Remove:
			_, span := p.tracer.Start(ctx, "pool remove")

			delete(p.Connections, fd)

			p.logger.Info("removed client from pool", slog.Int("connections", len(p.Connections)))
			span.End()
		case event := <-p.EpollEvents:
			_, span := p.tracer.Start(ctx, "pool epollevents")

			connections := make([]*Client, 0, len(event.Events))
			for i := range event.Events {
				conn, ok := p.Connections[event.Events[i].Fd]
				if !ok {
					continue
				}

				connections = append(connections, conn)
			}

			event.Resp <- connections
			span.End()
		case req := <-p.Snapshot:
			_, span := p.tracer.Start(ctx, "pool snapshot")

			s := &Snapshot{
				Connections: p.Connections,
				Processed:   make(chan struct{}),
			}

			req <- s
			<-s.Processed

			span.End()
		case event := <-p.Broadcast:
			ctx, span := p.tracer.Start(event.Ctx, "pool broadcast")

			p.logger.Info("received event to broadcast", slog.String("msg", event.Msg.Value))

			sendMsg := func(ctx context.Context, client *Client) error {
				return client.SendMessage(ctx, event.Msg)
			}

			stream := pipelines.StreamMap[int32, *Client](ctx, p.Connections)
			fanOut := pipelines.FanOut(ctx, stream, sendMsg, MaxSendFan)
			errChan := pipelines.FanIn(ctx, fanOut...)

			for err := range errChan {
				select {
				case <-ctx.Done():
					return
				default:
					if err != nil {
						p.logger.Error("failed to send message to client", slog.String("err", err.Error()))
					}
				}
			}

			event.Consumed <- struct{}{}
			span.End()
		}
	}
}
