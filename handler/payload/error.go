package payload

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
