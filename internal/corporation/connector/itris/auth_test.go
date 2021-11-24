package itris_test

import (
	"testing"

	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector/itris"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/julienrbrt/woningfinder/pkg/mapbox"
	"github.com/stretchr/testify/assert"
)

func Test_Login_Failed(t *testing.T) {
	a := assert.New(t)
	mockMapbox := mapbox.NewClientMock(nil, "district")

	client, err := itris.NewClient(logging.NewZapLoggerWithoutSentry(), mockMapbox, itris.OnsHuisInfo)
	a.NoError(err)

	err = client.Login("example", "unexesting")
	a.Error(connector.ErrAuthFailed, err)
}
