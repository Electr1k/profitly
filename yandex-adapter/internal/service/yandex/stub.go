package yandex

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

const geocodeStubPath = "internal/service/yandex/testdata/geocode.json"

type Stub struct {
	geocode GeocodeResult
}

func NewStub() (*Stub, error) {
	data, err := os.ReadFile(geocodeStubPath)
	if err != nil {
		return nil, fmt.Errorf("read geocode stub: %w", err)
	}

	var resp GeocodeResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse geocode stub: %w", err)
	}

	return &Stub{geocode: toGeocodeResult(&resp)}, nil
}

func (s *Stub) Geocode(_ context.Context, _ GeocodeParams) (GeocodeResult, error) {
	return s.geocode, nil
}
