package entity

// the structs here defines the refer table

type CorporationCity struct {
	CorporationName string
	CityName        string
}

type HousingPreferencesHousingType struct {
	HousingPreferencesID uint
	HousingType          string
}

type HousingPreferencesCityDistrict struct {
	HousingPreferencesID uint
	Name                 string
	CityName             string
}

type HousingPreferencesCity struct {
	HousingPreferencesID uint
	CityName             string
}
