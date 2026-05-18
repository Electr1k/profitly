package yandex

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

const (
	GeocodePath = "testdata/geocode.json"
)

type Stub struct {
	geocode GeocodeResult
}

func NewStub() (*Stub, error) {
	data, err := os.ReadFile(GeocodePath)
	if err != nil {
		return nil, fmt.Errorf("error on read geocode %w", err)
	}

	var resp GeocodeResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("error on parse geocode %w", err)
	}

	return &Stub{geocode: toGeocodeResult(&resp)}, nil
}

func (s *Stub) Geocode(_ context.Context, _ GeocodeParams) (GeocodeResult, error) {
	return s.geocode, nil
}
