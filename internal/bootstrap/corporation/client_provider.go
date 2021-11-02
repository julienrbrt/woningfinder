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
			Corporation: woningnet.AlmereInfo,
			ClientFunc: func() connector.Client {
				return CreateWoningNetClient(logger, mapboxClient, woningnet.AlmereInfo)
			},
		},
		{
			Corporation: woningnet.WoonkeusInfo,
			ClientFunc: func() connector.Client {
				return CreateWoningNetClient(logger, mapboxClient, woningnet.WoonkeusInfo)
			},
		},
		{
			Corporation: woningnet.EemvalleiInfo,
			ClientFunc: func() connector.Client {
				return CreateWoningNetClient(logger, mapboxClient, woningnet.EemvalleiInfo)
			},
		},
		{
			Corporation: woningnet.WoonserviceInfo,
			ClientFunc: func() connector.Client {
				return CreateWoningNetClient(logger, mapboxClient, woningnet.WoonserviceInfo)
			},
		},
		{
			Corporation: woningnet.MercatusInfo,
			ClientFunc: func() connector.Client {
				return CreateWoningNetClient(logger, mapboxClient, woningnet.MercatusInfo)
			},
		},
		{
			Corporation: woningnet.MiddenHollandInfo,
			ClientFunc: func() connector.Client {
				return CreateWoningNetClient(logger, mapboxClient, woningnet.MiddenHollandInfo)
			},
		},
		{
			Corporation: woningnet.BovenGroningenInfo,
			ClientFunc: func() connector.Client {
				return CreateWoningNetClient(logger, mapboxClient, woningnet.BovenGroningenInfo)
			},
		},
		{
			Corporation: woningnet.GooiVechtstreekInfo,
			ClientFunc: func() connector.Client {
				return CreateWoningNetClient(logger, mapboxClient, woningnet.GooiVechtstreekInfo)
			},
		},
		{
			Corporation: woningnet.GroningenInfo,
			ClientFunc: func() connector.Client {
				return CreateWoningNetClient(logger, mapboxClient, woningnet.GroningenInfo)
			},
		},
		{
			Corporation: woningnet.HuiswaartsInfo,
			ClientFunc: func() connector.Client {
				return CreateWoningNetClient(logger, mapboxClient, woningnet.HuiswaartsInfo)
			},
		},
		{
			Corporation: woningnet.WoongaardInfo,
			ClientFunc: func() connector.Client {
				return CreateWoningNetClient(logger, mapboxClient, woningnet.WoongaardInfo)
			},
		},
		{
			Corporation: zig.DeWoningZoekerInfo,
			ClientFunc:  func() connector.Client { return CreateZigClient(logger, mapboxClient, zig.DeWoningZoekerInfo) },
		},
		{
			Corporation: woonburo.AlmeloInfo,
			ClientFunc:  func() connector.Client { return CreateWoonburoClient(logger, mapboxClient, woonburo.AlmeloInfo) },
		},
		{
			Corporation: ikwilhuren.Info,
			ClientFunc:  func() connector.Client { return CreateIkWilHurenClient(logger, mapboxClient) },
		},
	}

	return connector.NewClientProvider(providers)
}
