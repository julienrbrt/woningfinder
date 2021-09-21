package city

import (
	"strings"
	"time"

	"github.com/woningfinder/woningfinder/pkg/logging"
)

// City defines a city where a housing corporation operates or when an house offer lies
// Note when crating an user city district are only contains in the key of the map
// This means that the matching only checks the key
// The values of the map is used for verbose suggestion of districts and for matching districts to the key if necessary
type City struct {
	CreatedAt time.Time           `pg:"default:now()" json:"-"`
	Name      string              `pg:",pk" json:"name"`
	District  map[string][]string `pg:"-" json:"district,omitempty"` // returns the abbrv district name and containing neighbourhood
}

func (c *City) Districts() []string {
	var districts []string
	for d := range c.District {
		districts = append(districts, d)
	}

	return districts
}

func (c *City) Neighbourhoods() []string {
	var neighbourhood []string
	for _, n := range c.District {
		neighbourhood = append(neighbourhood, n...)
	}

	return neighbourhood
}

// Merge cities that are supposed to be the same but that housing corporation name differently
func (c *City) Merge() City {
	switch {
	case strings.Contains(c.Name, "Hengelo"):
		return Hengelo
	case strings.Contains(c.Name, "Winterswijk"):
		return Winterswijk
	case strings.EqualFold(c.Name, "s-Heerenberg"):
		return Heerenberg
	}

	return *c
}

// cityTable defines the city name and it's corresponding struct
var cityTable = map[string]City{
	Aalten.Name:      Aalten,
	Dinxperlo.Name:   Dinxperlo,
	Bredevoort.Name:  Bredevoort,
	Neede.Name:       Neede,
	Borne.Name:       Borne,
	Bussum.Name:      Bussum,
	Wehl.Name:        Wehl,
	Enschede.Name:    Enschede,
	Haaksbergen.Name: Haaksbergen,
	Hengelo.Name:     Hengelo,
	Losser.Name:      Losser,
	Overdinkel.Name:  Overdinkel,
	DeLutte.Name:     DeLutte,
	Groenlo.Name:     Groenlo,
	Ulft.Name:        Ulft,
	Winterswijk.Name: Winterswijk,
	Zwolle.Name:      Zwolle,
	Almelo.Name:      Almelo,
	Hertme.Name:      Hertme,
	Apeldoorn.Name:   Apeldoorn,
	Heerenberg.Name:  Heerenberg,
	Zenderen.Name:    Zenderen,
	Utrecht.Name:     Utrecht,
	Zeist.Name:       Zeist,
}

// HasSuggestedCityDistrict permit to see if a city has an city districts
// This method is used for checking if we support a city districts before making the request to Mapbox (see in the connector logic).
func HasSuggestedCityDistrict(cityName string) bool {
	city, ok := cityTable[cityName]

	return ok && len(city.District) > 0
}

func SuggestedCityDistrictFromName(logger *logging.Logger, cityName string) map[string][]string {
	city, ok := cityTable[cityName]
	if !ok {
		logger.Sugar().Errorf("failed to get city district of %s", cityName)
		return nil
	}

	return city.District
}
