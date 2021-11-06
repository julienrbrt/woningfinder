package corporation

import (
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/dewoonplaats"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/domijn"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/ikwilhuren"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/itris"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/woningnet"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/woonburo"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/zig"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

// CreateClientProvider provides the client of a corporation
func CreateClientProvider(logger *logging.Logger, mapboxClient mapbox.Client) connector.ClientProvider {
	providers := []connector.Provider{
		{
			Corporation: dewoonplaats.Info,
			Connector:   CreateDeWoonplaatsClient(logger, mapboxClient),
		},
		{
			Corporation: itris.OnsHuisInfo,
			Connector:   CreateItrisClient(logger, mapboxClient, itris.OnsHuisInfo),
		},
		{
			Corporation: domijn.Info,
			Connector:   CreateDomijnClient(logger, mapboxClient),
		},
		{
			Corporation: zig.RoomspotInfo,
			Connector:   CreateZigClient(logger, mapboxClient, zig.RoomspotInfo),
		},
		{
			Corporation: zig.DeWoningZoekerInfo,
			Connector:   CreateZigClient(logger, mapboxClient, zig.DeWoningZoekerInfo),
		},
		{
			Corporation: zig.WoonnetHaaglanden,
			Connector:   CreateZigClient(logger, mapboxClient, zig.WoonnetHaaglanden),
		},
		{
			Corporation: woningnet.HengeloBorneInfo,
			Connector:   CreateWoningNetClient(logger, mapboxClient, woningnet.HengeloBorneInfo),
		},
		{
			Corporation: woningnet.UtrechtInfo,
			Connector:   CreateWoningNetClient(logger, mapboxClient, woningnet.UtrechtInfo),
		},
		{
			Corporation: woningnet.AmsterdamInfo,
			Connector:   CreateWoningNetClient(logger, mapboxClient, woningnet.AmsterdamInfo),
		},
		{
			Corporation: woningnet.AlmereInfo,
			Connector:   CreateWoningNetClient(logger, mapboxClient, woningnet.AlmereInfo),
		},
		{
			Corporation: woningnet.WoonkeusInfo,
			Connector:   CreateWoningNetClient(logger, mapboxClient, woningnet.WoonkeusInfo),
		},
		{
			Corporation: woningnet.EemvalleiInfo,
			Connector:   CreateWoningNetClient(logger, mapboxClient, woningnet.EemvalleiInfo),
		},
		{
			Corporation: woningnet.WoonserviceInfo,
			Connector:   CreateWoningNetClient(logger, mapboxClient, woningnet.WoonserviceInfo),
		},
		{
			Corporation: woningnet.MercatusInfo,
			Connector:   CreateWoningNetClient(logger, mapboxClient, woningnet.MercatusInfo),
		},
		{
			Corporation: woningnet.MiddenHollandInfo,
			Connector:   CreateWoningNetClient(logger, mapboxClient, woningnet.MiddenHollandInfo),
		},
		{
			Corporation: woningnet.BovenGroningenInfo,
			Connector:   CreateWoningNetClient(logger, mapboxClient, woningnet.BovenGroningenInfo),
		},
		{
			Corporation: woningnet.GooiVechtstreekInfo,
			Connector:   CreateWoningNetClient(logger, mapboxClient, woningnet.GooiVechtstreekInfo),
		},
		{
			Corporation: woningnet.GroningenInfo,
			Connector:   CreateWoningNetClient(logger, mapboxClient, woningnet.GroningenInfo),
		},
		{
			Corporation: woningnet.HuiswaartsInfo,
			Connector:   CreateWoningNetClient(logger, mapboxClient, woningnet.HuiswaartsInfo),
		},
		{
			Corporation: woningnet.WoongaardInfo,
			Connector:   CreateWoningNetClient(logger, mapboxClient, woningnet.WoongaardInfo),
		},
		{
			Corporation: woonburo.AlmeloInfo,
			Connector:   CreateWoonburoClient(logger, mapboxClient, woonburo.AlmeloInfo),
		},
		{
			Corporation: ikwilhuren.Info,
			Connector:   CreateIkWilHurenClient(logger, mapboxClient),
		},
	}

	return connector.NewClientProvider(providers)
}
