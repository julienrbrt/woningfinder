package corporation

import (
	"fmt"
	"log"
	"strings"

	"github.com/woningfinder/woningfinder/pkg/osm"
)

// DistrictName gets the district name from a location
func (o *Offer) DistrictName() string {
	if o.Housing.District != "" {
		return strings.ToLower(o.Housing.District)
	}

	var district string
	district, err := osm.GetResidential(fmt.Sprintf("%.5f", o.Housing.Location.Latitude), fmt.Sprintf("%.5f", o.Housing.Location.Longitude))
	if err != nil {
		log.Printf(fmt.Errorf("error getting district from %s: %w", o.Housing.Address, err).Error())
	}

	return district
}
