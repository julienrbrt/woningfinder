package corporation

import (
	"fmt"
	"net/url"
	"time"
)

// Corporation defines a housing corporations basic data
// That data is shared between every housing corporations
type Corporation struct {
	APIEndpoint     *url.URL  `pg:"-" json:",omitempty"`
	CreatedAt       time.Time `pg:"default:now()"`
	DeletedAt       time.Time `pg:",soft_delete" json:"-"`
	Name            string    `pg:",pk"`
	URL             string
	Cities          []City            `pg:"-"` // linked to CorporationCity
	SelectionMethod []SelectionMethod `pg:"-"`
	SelectionTime   time.Time
}

// HasMinimal ensure that the corporation contains the minimal required data
func (c *Corporation) HasMinimal() error {
	if c.Name == "" || c.URL == "" {
		return fmt.Errorf("corporation name or url missing")
	}

	if len(c.Cities) == 0 {
		return fmt.Errorf("corporation cities missing")
	}

	for _, city := range c.Cities {
		if city.Name == "" {
			return fmt.Errorf("corporation cities invalid")
		}
	}

	if _, err := url.Parse(c.URL); err != nil {
		return fmt.Errorf("corporation url invalid")
	}

	return nil
}

// City defines a city where a housing corporation operates or when an house offer lies
type City struct {
	CreatedAt time.Time      `pg:"default:now()" json:"-"`
	Name      string         `pg:",pk" json:"name"`
	District  []CityDistrict `pg:"rel:has-many,join_fk:city_name" json:"district,omitempty"`
}

type CityDistrict struct {
	CreatedAt time.Time `pg:"default:now()" json:"-"`
	CityName  string    `pg:",pk" json:"-"`
	Name      string    `pg:",pk" json:"name"`
}
