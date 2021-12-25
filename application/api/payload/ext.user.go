package payload

import "github.com/caravan-inc/fankey-server/domain/model/inventory"

func NewPostUserResponseFrom(user *inventory.User) PostUserResponse {
	return PostUserResponse{}
}
