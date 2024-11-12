package job

import (
	"time"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"github.com/julienrbrt/woningfinder/internal/corporation/scheduler"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// OffersChan is a channel to receive corporation offers
var OffersChan = make(chan corporation.Offers)

// HousingFinder populates the housing-finder cron jobs
func (j *Jobs) HousingFinder(c *cron.Cron, connectorProvider connector.ConnectorProvider) {
	// populate crons
	for _, corp := range connectorProvider.GetCorporations() {
		// get corporation client
		client, err := connectorProvider.GetClient(corp.Name)
		if err != nil {
			j.logger.Error("error while getting corporation client", zap.Error(err))
			continue
		}

		// schedule corporation offer fetching
		schedule := scheduler.CorporationScheduler(corp)
		for _, s := range schedule {
			c.Schedule(s, cron.FuncJob(func() {
				j.logger.Info("housing-finder job started", zap.String("corporation", corp.Name))

				ch := make(chan corporation.Offer)
				go func(ch chan corporation.Offer) {
					defer close(ch)

					if err := client.FetchOffers(ch); err != nil {
						j.logger.Error("error while fetching offers", zap.String("corporation", corp.Name), zap.Error(err))
					}
				}(ch)

				offers := corporation.Offers{
					CorporationName: corp.Name,
					Offer:           []corporation.Offer{},
				}

				// batch send offers every 5 seconds
				ticker := time.NewTicker(5 * time.Second)
				defer ticker.Stop()

				counter := 0
				for {
					select {
					case <-ticker.C:
						if len(offers.Offer) == 0 {
							continue
						}

						j.logger.Info("housing-finder job sending offers", zap.String("corporation", corp.Name), zap.Int("offers", len(offers.Offer)))
						OffersChan <- offers

						counter += len(offers.Offer)
						offers.Offer = []corporation.Offer{}
					case offer, ok := <-ch:
						if ok {
							offers.Offer = append(offers.Offer, offer)
						}

						if !ok && len(offers.Offer) == 0 {
							j.logger.Info("housing-finder job finished", zap.Int("offers sent", counter), zap.String("corporation", corp.Name))
							return
						}
					}
				}
			}))
		}
	}
}
