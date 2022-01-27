package payload

import "github.com/caravan-inc/fankey-server/domain/model/inventory"

func NewPostUserResponseFrom(_ *inventory.User) PostUserResponse {
	return PostUserResponse{}
}
