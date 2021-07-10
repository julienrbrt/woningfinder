package city

import (
	"strings"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// CityTable defines the city name and it's corresponding struct
var CityTable = map[string]corporation.City{
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
}

// Merge cities that are supposed to be the same but that housing corporation name differently
func Merge(city corporation.City) corporation.City {
	switch {
	case strings.Contains(city.Name, Hengelo.Name):
		return Hengelo
	case strings.Contains(city.Name, Winterswijk.Name):
		return Winterswijk
	}

	return city
}

// HasSuggestedCityDistrict permit to see if a city has an city districts
// This method is used for checking if we support a city distritcs before making the request to Mapbox (see in the connector logic).
func HasSuggestedCityDistrict(cityName string) bool {
	city, ok := CityTable[cityName]

	return ok && len(city.District) > 0
}

func SuggestedCityDistrictFromName(logger *logging.Logger, cityName string) []string {
	city, ok := CityTable[cityName]
	if !ok {
		logger.Sugar().Errorf("failed to get city district of %s", cityName)
		return nil
	}

	return city.District
}
