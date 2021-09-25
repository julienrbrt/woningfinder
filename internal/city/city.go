package city

import (
	"time"
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
