package domain

import (
	"context"

	"github.com/nxdir-s/IdleEngine/internal/core/entity"
	"github.com/nxdir-s/IdleEngine/internal/ports"
)

type Users struct {
	service ports.UserService
}

func NewUsers(service ports.UserService) *Users {
	return &Users{
		service: service,
	}
}

func (d *Users) GetUser(ctx context.Context, id int) (*entity.User, error) {
	return d.service.GetUser(ctx, id)
}
