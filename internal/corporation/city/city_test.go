package city_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
)

func TestCity_Merge(t *testing.T) {
	a := assert.New(t)

	a.Equal(((&city.City{Name: "Hengelo OV"}).Merge()), city.Hengelo)
	a.Equal(((&city.City{Name: "Hengelo"}).Merge()), city.Hengelo)
	expected := &city.City{Name: "a city"}
	a.Equal(*expected, expected.Merge())
}
