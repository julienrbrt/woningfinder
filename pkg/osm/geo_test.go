package osm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/pkg/osm"
)

func Test_GetResidential(t *testing.T) {
	a := assert.New(t)
	name, err := osm.GetResidential("52.23148", "6.89277")
	a.NoError(err)
	a.Equal("roombeek", name)
}
