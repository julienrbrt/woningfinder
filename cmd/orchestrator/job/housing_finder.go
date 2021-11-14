package job

import (
	"github.com/robfig/cron/v3"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/corporation/scheduler"
	"go.uber.org/zap"
)

// HousingFinder populates the housing-finder cron jobs
func (j *Jobs) HousingFinder(c *cron.Cron, clientProvider connector.ClientProvider) {
	// populate crons
	for _, corp := range clientProvider.List() {
		corp := corp // https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable

		// get corporation client
		client, err := clientProvider.Get(corp.Name)
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
					if err := client.FetchOffers(ch); err != nil {
						j.logger.Error("error while fetching offers", zap.String("corporation", corp.Name), zap.Error(err))
					}
					close(ch)
				}(ch)

				offers := corporation.Offers{
					CorporationName: corp.Name,
					Offer:           []corporation.Offer{},
				}

				// batch send 50 offers
				for offer := range ch {
					offers.Offer = append(offers.Offer, offer)

					if len(offers.Offer) == 50 {
						j.logger.Info("housing-finder job sending offers", zap.String("corporation", corp.Name), zap.Int("offers", len(offers.Offer)))

						if err := j.matcherService.SendOffers(offers); err != nil {
							j.logger.Error("error while sending offer to redis queue", zap.String("corporation", offer.CorporationName), zap.Error(err))
						}
						offers.Offer = []corporation.Offer{}
					}
				}
			}))
		}
	}
}
