package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(geocode *GeocodeHandler) chi.Router {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/geocode", geocode.Handle)
	})

	return r
}
