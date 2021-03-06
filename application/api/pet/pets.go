package pet

import (
	"net/http"

	"github.com/HMasataka/beyond/application/api/payload"
)

func NewPet() ServerInterface {
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
