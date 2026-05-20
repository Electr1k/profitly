package http

import (
	"encoding/json"
	"fias-adapter/internal/service/fias"
	"fias-adapter/internal/usecase/addressDetail"
	"log/slog"
	"net/http"
)

type AddressDetailInterface interface {
	Handle(input addressDetail.Input) (fias.Home, error)
}

type AddressDetailHandler struct {
	uc  AddressDetailInterface
	log *slog.Logger
}

func NewAddressDetailHandler(uc AddressDetailInterface, log *slog.Logger) *AddressDetailHandler {
	return &AddressDetailHandler{uc, log}
}

func (h *AddressDetailHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req addressDetail.Input
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.uc.Handle(req)
	if err != nil {
		h.log.Error("addressDetail usecase failed", "error", err)
		writeError(w, http.StatusBadGateway, "upstream addressDetail error")
		return
	}

	writeJSON(w, http.StatusOK, res)

}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
