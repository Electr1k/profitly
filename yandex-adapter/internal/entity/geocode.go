package entity

type Kind string

const (
	KindHouse    Kind = "house"
	KindStreet   Kind = "street"
	KindMetro    Kind = "metro"
	KindDistrict Kind = "district"
	KindLocality Kind = "locality"
)

func (k Kind) Valid() bool {
	switch k {
	case KindHouse, KindStreet, KindMetro, KindDistrict, KindLocality:
		return true
	}
	return false
}

type GeocodeRequest struct {
	Lng    float64
	Lat    float64
	Limit  int
	Offset int
	Kind   Kind
}

type GeoObject struct {
	Name      string
	Kind      string
	Precision string
	Address   string

	Country  *string
	Province *string
	Area     *string
	Locality *string
	District *string
	Street   *string
	House    *string
	Other    *string
}

type GeocodeResult struct {
	Items  []GeoObject
	Total  int
	Count  int
	Offset int
}
