package job

import (
	"github.com/robfig/cron/v3"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/corporation/scheduler"
)

// HousingFinder populates the housing-finder cron jobs
func (j *Jobs) HousingFinder(c *cron.Cron, clientProvider connector.ClientProvider) {
	// populate crons
	for _, corp := range clientProvider.List() {
		corp := corp // https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable

		// get corporation client
		client, err := clientProvider.Get(corp.Name)
		if err != nil {
			j.logger.Sugar().Error(err)
			continue
		}

		// schedule corporation offer fetching
		schedule := scheduler.CorporationScheduler(corp)
		for _, s := range schedule {
			c.Schedule(s, cron.FuncJob(func() {
				j.logger.Sugar().Infof("housing-finder '%s' job started", corp.Name)

				if err := j.matcherService.PushOffers(client, corp); err != nil {
					j.logger.Sugar().Error(err)
				}
			}))
		}
	}
}
