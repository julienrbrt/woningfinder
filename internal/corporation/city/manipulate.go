package city

import (
	"fmt"
	"strings"

	"github.com/woningfinder/woningfinder/internal/corporation"
)

var cityDistrictTable = map[string]corporation.City{
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
	Groenlo.Name:     Groenlo,
	Ulft.Name:        Ulft,
	Winterswijk.Name: Winterswijk,
	Zwolle.Name:      Zwolle,
}

// Merge cities that are supposed to be the same but that housing corporation name differently
func Merge(city corporation.City) corporation.City {
	switch {
	case strings.Contains(city.Name, "Hengelo"):
		return Hengelo
	}

	return city
}

func SuggestedCityDistrictFromName(cityName string) ([]string, error) {
	city, ok := cityDistrictTable[cityName]
	if !ok {
		return nil, fmt.Errorf("failed to get city district of %s", cityName)
	}

	return city.District, nil
}
