package domain

import (
	"context"
	"strconv"

	"github.com/nxdir-s/IdleEngine/internal/ports"
	"github.com/nxdir-s/IdleEngine/protobuf"
)

type ErrUnknownAction struct {
	action int
}

func (e *ErrUnknownAction) Error() string {
	return "unknown action: " + strconv.Itoa(e.action)
}

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
