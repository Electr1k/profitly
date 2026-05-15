package yandex

type GeocodeResponse struct {
	Response struct {
		GeoObjectCollection GeoObjectCollection `json:"GeoObjectCollection"`
	} `json:"response"`
}

type GeoObjectCollection struct {
	MetaDataProperty struct {
		GeocoderResponseMetaData ResponseMetaData `json:"GeocoderResponseMetaData"`
	} `json:"metaDataProperty"`
	FeatureMember []FeatureMember `json:"featureMember"`
}

type ResponseMetaData struct {
	Request string `json:"request"`
	Results string `json:"results"`
	Found   string `json:"found"`
}

type FeatureMember struct {
	GeoObject GeoObject `json:"GeoObject"`
}

type GeoObject struct {
	Name             string `json:"name"`
	MetaDataProperty struct {
		GeocoderMetaData GeocoderMetaData `json:"GeocoderMetaData"`
	} `json:"metaDataProperty"`
}

type GeocoderMetaData struct {
	Precision string  `json:"precision"`
	Text      string  `json:"text"`
	Kind      string  `json:"kind"`
	Address   Address `json:"Address"`
}

type Address struct {
	Formatted  string      `json:"formatted"`
	Components []Component `json:"Components"`
}

type Component struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}
