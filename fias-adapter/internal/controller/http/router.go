package http

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter(addressDetail *AddressDetailHandler) chi.Router {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Get("/address-detail", addressDetail.Handle)
	})

	return r
}
