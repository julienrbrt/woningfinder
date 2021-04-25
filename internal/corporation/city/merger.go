package city

import (
	"strings"

	"github.com/woningfinder/woningfinder/internal/corporation"
)

// Merge cities that are supposed to be the same but that housing corporation name differently
func Merge(city corporation.City) corporation.City {
	switch {
	case strings.Contains(city.Name, "Hengelo"):
		return Hengelo
	}

	return city
}
