package yandex

import (
	"strconv"

	"yandex-adapter/internal/entity"
)

func toGeocodeResult(resp *GeocodeResponse) GeocodeResult {
	members := resp.Response.GeoObjectCollection.FeatureMember
	items := make([]entity.GeoObject, 0, len(members))
	for _, m := range members {
		items = append(items, toGeoObject(m.GeoObject))
	}

	meta := resp.Response.GeoObjectCollection.MetaDataProperty.GeocoderResponseMetaData
	total, _ := strconv.Atoi(meta.Found)

	return GeocodeResult{
		Items: items,
		Total: total,
	}
}

func toGeoObject(g GeoObject) entity.GeoObject {
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
