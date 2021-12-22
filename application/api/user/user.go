package user

import (
	"fmt"
	"net/http"
)

func NewUser() *User {
	return &User{}
}

type User struct {
}

func (u *User) PostUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(string)
	w.Write([]byte(fmt.Sprintf("hi %s", user)))
}
