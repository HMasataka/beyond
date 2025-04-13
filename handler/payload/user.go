package payload

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
}

type GetUserResponse struct {
	ID string `json:"id"`
}

type PostUserRequest struct {
	Name string `json:"name"`
}

type PostUserResponse struct {
	ID string `json:"id"`
}
