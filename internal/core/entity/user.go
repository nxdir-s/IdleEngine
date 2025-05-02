package entity

import "github.com/nxdir-s/IdleEngine/internal/core/entity/users"

type User struct {
	Id    int32
	State *users.State
}

func NewUser() *User {
	return &User{
		State: users.NewState(),
	}
}
