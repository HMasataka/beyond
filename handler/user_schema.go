package handler

import (
	"github.com/HMasataka/beyond/domain/entity"
	"github.com/HMasataka/beyond/domain/inventory"
	"github.com/HMasataka/beyond/handler/payload"
)

func NewUserFrom(user *entity.User) payload.User {
	if user == nil {
		return payload.User{}
	}

	return payload.User{
		ID:        user.ID,
		Name:      user.Name,
		Icon:      user.Icon,
		CreatedAt: user.CreatedAt,
	}
}

func NewUser(user *inventory.User) payload.User {
	if user == nil {
		return payload.User{}
	}

	return payload.User{
		ID:        user.ID,
		Name:      user.Name,
		Icon:      user.Icon,
		CreatedAt: user.CreatedAt,
	}
}
