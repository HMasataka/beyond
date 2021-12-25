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

func (p *Pet) ListPets(w http.ResponseWriter, r *http.Request, params ListPetsParams) (interface{}, error) {
	id := "pet_id"
	return payload.ListPetResponse{Id: &id}, nil
}

func (p *Pet) FindPets(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, nil
}
