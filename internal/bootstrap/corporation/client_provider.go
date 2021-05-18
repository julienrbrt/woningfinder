package corporation

import (
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/dewoonplaats"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/domijn"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/itris"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/woningnet"
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
			Client:      CreateItrisClient(logger, mapboxClient, itris.OnsHuisInfo),
		},
		{
			Corporation: zig.RoomspotInfo,
			Client:      CreateZigClient(logger, mapboxClient, zig.RoomspotInfo),
		},
		{
			Corporation: domijn.Info,
			Client:      CreateDomijnClient(logger, mapboxClient),
		},
		{
			Corporation: woningnet.HengeloBorneInfo,
			Client:      CreateWoningNetClient(logger, mapboxClient, woningnet.HengeloBorneInfo),
		},
	}

	return connector.NewClientProvider(providers)
}
