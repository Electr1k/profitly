package geocode

import (
	"context"
	"fmt"
	"strconv"

	"yandex-adapter/internal/entity"
	"yandex-adapter/internal/service/yandex"
)

type YandexClient interface {
	Geocode(ctx context.Context, p yandex.GeocodeParams) (*yandex.GeocodeResponse, error)
}

type UseCase struct {
	client YandexClient
}

func New(client YandexClient) *UseCase {
	return &UseCase{client: client}
}

func (u *UseCase) Get(ctx context.Context, req entity.GeocodeRequest) (entity.GeocodeResult, error) {
	resp, err := u.client.Geocode(ctx, yandex.GeocodeParams{
		Lng:     req.Lng,
		Lat:     req.Lat,
		Results: req.Limit,
		Skip:    req.Offset,
		Kind:    string(req.Kind),
	})
	if err != nil {
		return entity.GeocodeResult{}, fmt.Errorf("yandex geocode: %w", err)
	}

	members := resp.Response.GeoObjectCollection.FeatureMember
	items := make([]entity.GeoObject, 0, len(members))
	for _, m := range members {
		items = append(items, toEntity(m.GeoObject))
	}

	meta := resp.Response.GeoObjectCollection.MetaDataProperty.GeocoderResponseMetaData
	total, _ := strconv.Atoi(meta.Found)

	return entity.GeocodeResult{
		Items:  items,
		Total:  total,
		Count:  len(items),
		Offset: req.Offset,
	}, nil
}

func toEntity(g yandex.GeoObject) entity.GeoObject {
	md := g.MetaDataProperty.GeocoderMetaData
	obj := entity.GeoObject{
		Name:      g.Name,
		Kind:      md.Kind,
		Precision: md.Precision,
		Address:   md.Address.Formatted,
	}
	for _, c := range md.Address.Components {
		name := c.Name
		switch c.Kind {
		case "country":
			obj.Country = &name
		case "province":
			obj.Province = &name
		case "area":
			obj.Area = &name
		case "locality":
			obj.Locality = &name
		case "district":
			obj.District = &name
		case "street":
			obj.Street = &name
		case "house":
			obj.House = &name
		case "other":
			obj.Other = &name
		}
	}
	return obj
}
