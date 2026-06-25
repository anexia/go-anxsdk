package ipam

// Location represents a location in ipam responses.
type Location struct {
	ID         int    `json:"id"`
	Identifier string `json:"identifier"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Country    string `json:"country"`
	Lat        string `json:"lat"`
	Lon        string `json:"lon"`
	CityCode   string `json:"city_code"`
	IsPhysical bool   `json:"is_physical"`
	NameFull   string `json:"name_full"`
}
