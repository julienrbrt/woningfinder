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
			ClientFunc:  func() connector.Client { return CreateDeWoonplaatsClient(logger, mapboxClient) },
		},
		{
			Corporation: itris.OnsHuisInfo,
			ClientFunc:  func() connector.Client { return CreateItrisClient(logger, mapboxClient, itris.OnsHuisInfo) },
		},
		{
			Corporation: zig.RoomspotInfo,
			ClientFunc:  func() connector.Client { return CreateZigClient(logger, mapboxClient, zig.RoomspotInfo) },
		},
		{
			Corporation: domijn.Info,
			ClientFunc:  func() connector.Client { return CreateDomijnClient(logger, mapboxClient) },
		},
		{
			Corporation: woningnet.HengeloBorneInfo,
			ClientFunc: func() connector.Client {
				return CreateWoningNetClient(logger, mapboxClient, woningnet.HengeloBorneInfo)
			},
		},
		{
			Corporation: woningnet.UtrechtInfo,
			ClientFunc: func() connector.Client {
				return CreateWoningNetClient(logger, mapboxClient, woningnet.UtrechtInfo)
			},
		},
		{
			Corporation: woningnet.AmsterdamInfo,
			ClientFunc: func() connector.Client {
				return CreateWoningNetClient(logger, mapboxClient, woningnet.AmsterdamInfo)
			},
		},
		{
			Corporation: zig.DeWoningZoekerInfo,
			ClientFunc:  func() connector.Client { return CreateZigClient(logger, mapboxClient, zig.DeWoningZoekerInfo) },
		},
	}

	return connector.NewClientProvider(providers)
}
