package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"yandex-adapter/internal/entity"
	"yandex-adapter/internal/usecase/geocode"
)

const (
	defaultLimit = 100
	maxLimit     = 100
)

type GeocodeUseCase interface {
	Get(ctx context.Context, input geocode.Input) (geocode.Output, error)
}

type GeocodeHandler struct {
	uc  GeocodeUseCase
	log *slog.Logger
}

func NewGeocodeHandler(uc GeocodeUseCase, log *slog.Logger) *GeocodeHandler {
	return &GeocodeHandler{uc: uc, log: log}
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
	req, err := parseRequest(r.URL.Query())
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

func parseRequest(q url.Values) (geocode.Input, error) {
	lngStr := q.Get("lng")
	latStr := q.Get("lat")
	if lngStr == "" || latStr == "" {
		return geocode.Input{}, errors.New("lng and lat are required")
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		return geocode.Input{}, fmt.Errorf("lng: invalid float: %w", err)
	}
	if lng < -180 || lng > 180 {
		return geocode.Input{}, errors.New("lng must be in [-180, 180]")
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return geocode.Input{}, fmt.Errorf("lat: invalid float: %w", err)
	}
	if lat < -90 || lat > 90 {
		return geocode.Input{}, errors.New("lat must be in [-90, 90]")
	}

	limit := defaultLimit
	if s := q.Get("limit"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			return geocode.Input{}, fmt.Errorf("limit: invalid int: %w", err)
		}
		if v < 1 {
			return geocode.Input{}, errors.New("limit must be >= 1")
		}
		limit = v
		if limit > maxLimit {
			limit = maxLimit
		}
	}

	offset := 0
	if s := q.Get("offset"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			return geocode.Input{}, fmt.Errorf("offset: invalid int: %w", err)
		}
		if v < 0 {
			return geocode.Input{}, errors.New("offset must be >= 0")
		}
		offset = v
	}

	var kind entity.Kind
	if s := q.Get("kind"); s != "" {
		kind = entity.Kind(s)
		if !kind.Valid() {
			return geocode.Input{}, errors.New("kind must be one of: house, street, metro, district, locality")
		}
	}

	return geocode.Input{
		Lng:    lng,
		Lat:    lat,
		Limit:  limit,
		Offset: offset,
		Kind:   kind,
	}, nil
}

func toResponseDTO(res geocode.Output) geocodeResponseDTO {
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
