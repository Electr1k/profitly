package geocode

import "yandex-adapter/internal/entity"

type Input struct {
	Lng    float64
	Lat    float64
	Limit  int
	Offset int
	Kind   entity.Kind
}

type Output struct {
	Items  []entity.GeoObject
	Total  int
	Count  int
	Offset int
}
