package payload

type PostAccountRequest struct{}

type PostAccountResponse struct {
	UserID string `json:"user_id"`
}
