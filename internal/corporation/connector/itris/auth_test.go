package itris_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/itris"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/itris/onshuis"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

func Test_Login(t *testing.T) {
	a := assert.New(t)
	mockMapbox := mapbox.NewClientMock(nil, "district")

	client, err := itris.NewClient(logging.NewZapLoggerWithoutSentry(), mockMapbox, "https://mijn.onshuis.com", onshuis.DetailsParser)
	a.NoError(err)

	a.NoError(client.Login("julien_", "de!"))
}
