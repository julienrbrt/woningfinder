package entity

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
	UpdatedAt       time.Time
	DeletedAt       time.Time `pg:",soft_delete"`
	Name            string    `pg:",pk"`
	URL             string
	Cities          []City            `pg:"many2many:corporation_cities"`
	SelectionMethod []SelectionMethod `pg:"many2many:corporation_selection_methods,join_fk:selection_method"`
	SelectionTime   time.Time
}

// IsValid verifies if the given corporation is valid
func (c *Corporation) IsValid() error {
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

	if len(c.SelectionMethod) == 0 {
		return fmt.Errorf("corporation selection method missing")
	}

	for _, selection := range c.SelectionMethod {
		if !selection.Method.Exists() {
			return fmt.Errorf("corporation selection method invalid")
		}
	}

	if _, err := url.Parse(c.URL); err != nil {
		return fmt.Errorf("corporation url invalid")
	}

	return nil
}

// City defines a city where a housing corporation operates or when an house offer lies
type City struct {
	CreatedAt time.Time `pg:"default:now()"`
	DeletedAt time.Time `pg:",soft_delete"`
	Name      string    `pg:",pk"`
	District  []string  `pg:"-" json:",omitempty"`
}

// CorporationCity defines the many-to-many relationship table
type CorporationCity struct {
	CorporationName string
	CityName        string
}

// CorporationSelectionMethod defines the many-to-many relationship table
type CorporationSelectionMethod struct {
	CorporationName string
	SelectionMethod Method
}
