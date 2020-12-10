package corporation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation"
)

func getOffer() corporation.Offer {
	return corporation.Offer{
		ExternalID: "w758752",
		Housing: corporation.Housing{
			Type:                    corporation.House,
			Latitude:                52.133,
			Longitude:               6.61433,
			Address:                 "Beatrixstraat 1 R 7161 DJ Neede  A",
			EnergieLabel:            "A",
			Price:                   656.39,
			Size:                    80,
			NumberRoom:              6,
			NumberBedroom:           2,
			BuildingYear:            2010,
			HousingAllowance:        true,
			Garden:                  false,
			Garage:                  false,
			Elevator:                true,
			Balcony:                 true,
			AccessibilityWheelchair: false,
			AccessibilityScooter:    true,
			Attic:                   false,
			Historic:                false,
			CV:                      false,
		},
		CanApply:   true,
		HasApplied: false,
	}
}

func Test_Offer_DistrictName(t *testing.T) {
	a := assert.New(t)
	offer := getOffer()
	offer.Housing.District = "roombeek"
	a.Equal("roombeek", offer.DistrictName())
}

func Test_Offer_DistrictName_FromOSM(t *testing.T) {
	a := assert.New(t)
	offer := getOffer()
	a.Equal("", offer.DistrictName())
	offer.Housing.Latitude = 52.23148
	offer.Housing.Longitude = 6.89277
	a.Equal("roombeek", offer.DistrictName())
}
