package pet

import (
	"net/http"

	"github.com/caravan-inc/fankey-server/application/api/payload"
)

func NewPet() *Pet {
	return &Pet{}
}

type Pet struct {
}

func (*Pet) ListPets(_ http.ResponseWriter, _ *http.Request, _ ListPetsParams) (interface{}, error) {
	id := "pet_id"
	return payload.ListPetResponse{Id: &id}, nil
}

func (*Pet) PostPet(_ http.ResponseWriter, _ *http.Request) (interface{}, error) {
	id := "pet_id"
	return payload.PostPetResponse{Id: &id}, nil
}
