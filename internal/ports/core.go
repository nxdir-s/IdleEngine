package ports

import (
	"context"

	"github.com/nxdir-s/IdleEngine/internal/core/entity"
	"github.com/nxdir-s/IdleEngine/protobuf"
)

type Events interface {
	HandleUserEvent(ctx context.Context, event *protobuf.UserEvent) error
}

type Users interface {
	GetUser(ctx context.Context, id int) (*entity.User, error)
}

type UserService interface {
	GetUser(ctx context.Context, id int) (*entity.User, error)
	GetUserID(ctx context.Context, email string) (int, error)
}

type UserTxService interface {
	UserService
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
