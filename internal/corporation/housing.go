package corporation

// Housing defines an appartement and a house
type Housing struct {
	Type          HousingType
	Address       string
	City          City
	CityDistrict  string
	EnergieLabel  string
	Price         float64
	Size          float64
	NumberRoom    int
	NumberBedroom int
	BuildingYear  int
	Garden        bool
	Garage        bool
	Elevator      bool
	Balcony       bool
	Attic         bool
	Accessible    bool // Assessible defines if the house is accessible for handicapt people
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
