package user

import (
	"net/http"

	"github.com/HMasataka/beyond/application/api/payload"
)

func NewUser() ServerInterface {
	return &User{}
}

type User struct {
}

func (*User) PostUser(_ http.ResponseWriter, _ *http.Request) (interface{}, error) {
	id := "pet_id"
	return payload.PostUserResponse{Id: &id}, nil
}
