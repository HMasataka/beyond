package user

import (
	"net/http"

	"github.com/caravan-inc/fankey-server/application/api/payload"
)

func NewUser() *User {
	return &User{}
}

type User struct {
}

func (u *User) PostUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id := "pet_id"
	return payload.PostUserResponse{Id: &id}, nil
}
