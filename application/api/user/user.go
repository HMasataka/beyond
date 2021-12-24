package user

import (
	"net/http"
)

func NewUser() *User {
	return &User{}
}

type User struct {
}

func (u *User) PostUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, nil
}
