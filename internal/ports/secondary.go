package ports

import (
	"context"

	"github.com/nxdir-s/IdleEngine/internal/adapters/secondary/franz"
	"github.com/nxdir-s/IdleEngine/internal/core/entity"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Kafka interface {
	Send(ctx context.Context, record protoreflect.ProtoMessage) error
	Consume(ctx context.Context, consumer franz.Consumer)
	Close() error
}

type Database interface {
	NewTransactionAdapter(ctx context.Context) (DatabaseTx, error)
	GetUser(ctx context.Context, id int) (*entity.User, error)
	GetUserID(ctx context.Context, email string) (int, error)
}

type DatabaseTx interface {
	Database
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
