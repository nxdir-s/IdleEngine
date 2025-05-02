package ports

import (
	"context"

	"github.com/nxdir-s/IdleEngine/protobuf"
)

type Consumer interface {
	ProcessUserEvent(ctx context.Context, event *protobuf.UserEvent) error
}
