package geocode

import (
	"context"
	"fmt"

	"yandex-adapter/internal/service/yandex"
)

type YandexClient interface {
	Geocode(ctx context.Context, p yandex.GeocodeParams) (yandex.GeocodeResult, error)
}

type UseCase struct {
	client YandexClient
}

func New(client YandexClient) *UseCase {
	return &UseCase{client: client}
}

func (u *UseCase) Get(ctx context.Context, input Input) (Output, error) {
	res, err := u.client.Geocode(ctx, yandex.GeocodeParams{
		Lng:     input.Lng,
		Lat:     input.Lat,
		Results: input.Limit,
		Skip:    input.Offset,
		Kind:    string(input.Kind),
	})
	if err != nil {
		return Output{}, fmt.Errorf("yandex geocode: %w", err)
	}

	return Output{
		Items:  res.Items,
		Total:  res.Total,
		Count:  len(res.Items),
		Offset: input.Offset,
	}, nil
}
