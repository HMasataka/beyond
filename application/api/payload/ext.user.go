package payload

import "github.com/HMasataka/beyond/domain/model/inventory"

func NewPostUserResponseFrom(_ *inventory.User) PostUserResponse {
	return PostUserResponse{}
}
