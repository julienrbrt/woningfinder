package customer

// the structs here defines the refer table
type HousingPreferencesHousingType struct {
	UserID      uint
	HousingType string
}

type HousingPreferencesCityDistrict struct {
	UserID   uint
	CityName string
	Name     string
}

type HousingPreferencesCity struct {
	UserID   uint
	CityName string
}
