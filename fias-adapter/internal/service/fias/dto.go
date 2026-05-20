package fias

type Address struct {
	ObjectId      int    `json:"object_id"`
	Path          string `json:"path"`
	ObjectLevelId int    `json:"object_level_id"`
	ObjectGuid    string `json:"object_guid"`
	FullName      string `json:"full_name"`
	AddressType   int    `json:"address_type"`
	Hierarchy     []struct {
		Name          string `json:"name"`
		FullName      string `json:"full_name"`
		Number        string `json:"number"`
		ObjectLevelId int    `json:"object_level_id"`
		ObjectId      int    `json:"object_id"`
	} `json:"hierarchy"`
}

type GetAddressItemsResponse struct {
	Addresses []Address `json:"addresses"`
}

type Home struct {
	Name          string `json:"name"`
	ApparentCount int    `json:"apparent_count"`
	ParkingCount  int    `json:"parking_count"`
}
