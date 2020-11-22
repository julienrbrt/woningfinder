package corporation

const (
	House HousingType = iota
	Appartement
	Parking
)

// HousingType defines the type of an Housing (appartement, house,)
type HousingType int

// Housing defines an appartement and a house
type Housing struct {
	Address       string
	Location      Location
	Type          HousingType
	Facilities    Facilities
	EnergieLabel  string
	Room, Bedroom int
	Price         int
	// Allowance defines if the house can get huurtoeslag
	Allowance bool
}

// Facilities are extra characteristics that belongs to a Housing
type Facilities struct {
	Garden   bool
	Garage   bool
	Terasse  bool
	Elevator bool
	Balcony  bool
}
