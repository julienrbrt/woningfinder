package corporation

import (
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

// CreateClientProvider provides the client of a corporation
func CreateClientProvider(logger *logging.Logger, mapboxClient mapbox.Client) connector.ClientProvider {
	providers := []connector.Provider{
		{
			Corporation: DeWoonplaatsInfo,
			Client:      CreateDeWoonplaatsClient(logger, mapboxClient),
		},
		{
			Corporation: OnsHuisInfo,
			Client:      CreateOnsHuisClient(logger, mapboxClient),
		},
		{
			Corporation: RoomspotInfo,
			Client:      CreateRoomspotClient(logger, mapboxClient),
		},
	}

	return connector.NewClientProvider(providers)
}
