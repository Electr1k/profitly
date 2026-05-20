package addressDetail

type Input struct {
	Region        string `json:"region"`
	AddressObject string `json:"address_object"`
	City          string `json:"city"`
	Street        string `json:"street"`
	Home          string `json:"home"`
}

type Output struct {
	ObjectId      int    `json:"object_id"`
	Path          string `json:"path"`
	ObjectLevelId int    `json:"object_level_id"`
	ObjectGuid    string `json:"object_guid"`
	FullName      string `json:"full_name"`
	AddressType   int    `json:"address_type"`
}
