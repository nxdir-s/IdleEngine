package service

import (
	"context"

	"github.com/nxdir-s/IdleEngine/internal/core/entity"
	"github.com/nxdir-s/IdleEngine/internal/ports"
)

type UserServiceOpt func(s *UserService)

func WithUserTx(tx ports.DatabaseTx) UserServiceOpt {
	return func(s *UserService) {
		s.tx = tx
	}
}

type UserService struct {
	db ports.Database
	tx ports.DatabaseTx
}

func NewUserService(database ports.Database, opts ...UserServiceOpt) *UserService {
	service := &UserService{
		db: database,
	}

	for _, opt := range opts {
		opt(service)
	}

	return service
}

func (s *UserService) NewUserTxService(ctx context.Context) (ports.UserTxService, error) {
	txAdapter, err := s.db.NewTransactionAdapter(ctx)
	if err != nil {
		return nil, err
	}

	txService := NewUserService(txAdapter, WithUserTx(txAdapter))

	return txService, nil
}

func (s *UserService) Commit(ctx context.Context) error {
	return s.tx.Commit(ctx)
}

func (s *UserService) Rollback(ctx context.Context) error {
	return s.tx.Rollback(ctx)
}

func (s *UserService) GetUser(ctx context.Context, id int) (*entity.User, error) {
	return s.db.GetUser(ctx, id)
}

func (s *UserService) GetUserID(ctx context.Context, email string) (int, error) {
	return s.db.GetUserID(ctx, email)
}
