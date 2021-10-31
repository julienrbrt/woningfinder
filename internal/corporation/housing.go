package corporation

// Housing defines an appartement and a house
type Housing struct {
	Type          HousingType
	Address       string
	CityName      string
	CityDistrict  string
	Price         float64
	Size          float64
	NumberBedroom int
	BuildingYear  int
	Garden        bool
	Garage        bool
	Elevator      bool
	Balcony       bool
	Accessible    bool // Accessible defines if the house is accessible for handicapt people
}

// HousingType defines the type of an HousingType (appartement, house)
type HousingType string

const (
	// HousingTypeHouse is a house type
	HousingTypeHouse HousingType = "house"
	// HousingTypeAppartement is a appartement type
	HousingTypeAppartement HousingType = "appartement"
	// HousingTypeUndefined is an undefined type (probably parking but unsupported by WoningFinder)
	HousingTypeUndefined HousingType = "undefined"
)
