package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"yandex-adapter/internal/entity"
)

const (
	defaultLimit = 100
	maxLimit     = 100
)

type GeocodeUseCase interface {
	Get(ctx context.Context, req entity.GeocodeRequest) (entity.GeocodeResult, error)
}

type GeocodeHandler struct {
	uc  GeocodeUseCase
	log *slog.Logger
}

func NewGeocodeHandler(uc GeocodeUseCase, log *slog.Logger) *GeocodeHandler {
	return &GeocodeHandler{uc: uc, log: log}
}

type geocodeRequestDTO struct {
	Lng    *float64 `json:"lng"`
	Lat    *float64 `json:"lat"`
	Limit  *int     `json:"limit"`
	Offset *int     `json:"offset"`
	Kind   *string  `json:"kind"`
}

type geoObjectDTO struct {
	Name      string  `json:"name"`
	Kind      string  `json:"kind"`
	Precision string  `json:"precision"`
	Address   string  `json:"address"`
	Country   *string `json:"country"`
	Province  *string `json:"province"`
	Area      *string `json:"area"`
	Locality  *string `json:"locality"`
	District  *string `json:"district"`
	Street    *string `json:"street"`
	House     *string `json:"house"`
	Other     *string `json:"other"`
}

type paginationDTO struct {
	Total  int `json:"total"`
	Count  int `json:"count"`
	Offset int `json:"offset"`
}

type geocodeResponseDTO struct {
	Items      []geoObjectDTO `json:"items"`
	Pagination paginationDTO  `json:"pagination"`
}

func (h *GeocodeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var body geocodeRequestDTO
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
		return
	}

	req, err := parseRequest(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.uc.Get(r.Context(), req)
	if err != nil {
		h.log.Error("geocode usecase failed", "error", err)
		writeError(w, http.StatusBadGateway, "upstream geocoder error")
		return
	}

	writeJSON(w, http.StatusOK, toResponseDTO(res))
}

func parseRequest(b geocodeRequestDTO) (entity.GeocodeRequest, error) {
	if b.Lng == nil || b.Lat == nil {
		return entity.GeocodeRequest{}, errors.New("lng and lat are required")
	}
	if *b.Lng < -180 || *b.Lng > 180 {
		return entity.GeocodeRequest{}, errors.New("lng must be in [-180, 180]")
	}
	if *b.Lat < -90 || *b.Lat > 90 {
		return entity.GeocodeRequest{}, errors.New("lat must be in [-90, 90]")
	}

	limit := defaultLimit
	if b.Limit != nil {
		if *b.Limit < 1 {
			return entity.GeocodeRequest{}, errors.New("limit must be >= 1")
		}
		limit = *b.Limit
		if limit > maxLimit {
			limit = maxLimit
		}
	}

	offset := 0
	if b.Offset != nil {
		if *b.Offset < 0 {
			return entity.GeocodeRequest{}, errors.New("offset must be >= 0")
		}
		offset = *b.Offset
	}

	var kind entity.Kind
	if b.Kind != nil && *b.Kind != "" {
		kind = entity.Kind(*b.Kind)
		if !kind.Valid() {
			return entity.GeocodeRequest{}, errors.New("kind must be one of: house, street, metro, district, locality")
		}
	}

	return entity.GeocodeRequest{
		Lng:    *b.Lng,
		Lat:    *b.Lat,
		Limit:  limit,
		Offset: offset,
		Kind:   kind,
	}, nil
}

func toResponseDTO(res entity.GeocodeResult) geocodeResponseDTO {
	items := make([]geoObjectDTO, 0, len(res.Items))
	for _, o := range res.Items {
		items = append(items, geoObjectDTO{
			Name:      o.Name,
			Kind:      o.Kind,
			Precision: o.Precision,
			Address:   o.Address,
			Country:   o.Country,
			Province:  o.Province,
			Area:      o.Area,
			Locality:  o.Locality,
			District:  o.District,
			Street:    o.Street,
			House:     o.House,
			Other:     o.Other,
		})
	}
	return geocodeResponseDTO{
		Items: items,
		Pagination: paginationDTO{
			Total:  res.Total,
			Count:  res.Count,
			Offset: res.Offset,
		},
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
