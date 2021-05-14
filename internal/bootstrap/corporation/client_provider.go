package corporation

import (
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/dewoonplaats"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/domijn"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/itris"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/zig"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

// CreateClientProvider provides the client of a corporation
func CreateClientProvider(logger *logging.Logger, mapboxClient mapbox.Client) connector.ClientProvider {
	providers := []connector.Provider{
		{
			Corporation: dewoonplaats.Info,
			Client:      CreateDeWoonplaatsClient(logger, mapboxClient),
		},
		{
			Corporation: itris.OnsHuisInfo,
			Client:      CreateOnsHuisClient(logger, mapboxClient),
		},
		{
			Corporation: zig.RoomspotInfo,
			Client:      CreateRoomspotClient(logger, mapboxClient),
		},
		{
			Corporation: domijn.Info,
			Client:      CreateDomijnClient(logger, mapboxClient),
		},
	}

	return connector.NewClientProvider(providers)
}
