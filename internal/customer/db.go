package customer

// the structs here defines the refer table
type HousingPreferencesHousingType struct {
	HousingPreferencesID uint
	HousingType          string
}

type HousingPreferencesCityDistrict struct {
	HousingPreferencesID uint
	CityName             string
	Name                 string
}

type HousingPreferencesCity struct {
	HousingPreferencesID uint
	CityName             string
}
