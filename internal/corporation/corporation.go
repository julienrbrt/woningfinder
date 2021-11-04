package corporation

import (
	"fmt"
	"net/url"
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation/city"
)

// Corporation defines a housing corporations basic data
// That data is shared between every housing corporations
type Corporation struct {
	APIEndpoint     *url.URL  `pg:"-" json:",omitempty"`
	CreatedAt       time.Time `pg:"default:now()"`
	DeletedAt       time.Time `pg:",soft_delete" json:"-"`
	Name            string    `pg:",pk"`
	URL             string
	Cities          []city.City       `pg:"-"` // linked to CorporationCity
	SelectionMethod []SelectionMethod `pg:"-"`
	SelectionTime   []time.Time       `pg:"-" json:"-"`
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
