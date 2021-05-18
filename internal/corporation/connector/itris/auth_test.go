package itris_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/itris"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

func Test_Login_Failed(t *testing.T) {
	a := assert.New(t)
	mockMapbox := mapbox.NewClientMock(nil, "district")

	client, err := itris.NewClient(logging.NewZapLoggerWithoutSentry(), mockMapbox, itris.OnsHuisInfo)
	a.NoError(err)

	err = client.Login("example", "unexesting")
	a.Error(connector.ErrAuthFailed, err)
}
