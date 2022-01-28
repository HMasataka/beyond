package inventory

import "github.com/HMasataka/beyond/domain/model"

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewUserFrom(m *model.User) User {
	return User{
		ID:   m.ID,
		Name: m.Name,
	}
}
