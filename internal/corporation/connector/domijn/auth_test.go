package domijn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/domijn"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

func Test_Login_Failed(t *testing.T) {
	a := assert.New(t)
	mockMapbox := mapbox.NewClientMock(nil, "district")

	client, err := domijn.NewClient(logging.NewZapLoggerWithoutSentry(), mockMapbox)
	a.NoError(err)

	err = client.Login("example@example.com", "unexesting")
	a.Error(connector.ErrAuthFailed, err)
}
