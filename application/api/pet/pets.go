package pet

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func NewPet() *Pet {
	return &Pet{}
}

type Pet struct {
}

func (p *Pet) ListPets(w http.ResponseWriter, r *http.Request, params ListPetsParams) {
	ctx := r.Context()

	user := ctx.Value("user").(string)
	requestID := middleware.GetReqID(ctx)

	w.Write([]byte(fmt.Sprintf("hi %s\n", user)))
	w.Write([]byte(fmt.Sprintf("your ID is %s\n", requestID)))
}

func (p *Pet) FindPets(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(string)
	w.Write([]byte(fmt.Sprintf("hi %s", user)))
}
