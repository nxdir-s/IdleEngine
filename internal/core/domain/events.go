package domain

import (
	"context"

	"github.com/nxdir-s/IdleEngine/internal/ports"
	"github.com/nxdir-s/IdleEngine/protobuf"
)

type Events struct {
	users ports.Users
}

func NewEvents(users ports.Users) *Events {
	return &Events{
		users: users,
	}
}

func (d *Events) HandleUserEvent(ctx context.Context, event *protobuf.UserEvent) error {
	return nil
}
