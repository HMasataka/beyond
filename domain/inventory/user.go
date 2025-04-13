package inventory

import (
	"time"

	"github.com/HMasataka/beyond/domain/entity"
)

func NewUser(user *entity.User) *User {
	if user == nil {
		return nil
	}

	return &User{
		ID:        user.ID,
		Name:      user.Name,
		Icon:      user.Icon,
		CreatedAt: user.CreatedAt,
	}
}

type User struct {
	ID        string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name      string    `boil:"name" json:"name" toml:"name" yaml:"name"`
	Icon      string    `boil:"icon" json:"icon" toml:"icon" yaml:"icon"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
}
