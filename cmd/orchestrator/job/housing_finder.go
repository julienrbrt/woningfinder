package job

import (
	"github.com/robfig/cron/v3"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/scheduler"
	matcherService "github.com/woningfinder/woningfinder/internal/services/matcher"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// HousingFinder populates the housing-finder cron jobs
func HousingFinder(logger *logging.Logger, c *cron.Cron, clientProvider corporation.ClientProvider, matcherService matcherService.Service) {
	// populate crons
	for _, corp := range clientProvider.List() {
		corp := corp // https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable

		// get corporation client
		client, err := clientProvider.GetByName(corp)
		if err != nil {
			logger.Sugar().Error(err)
			continue
		}

		// schedule corporation offer fetching
		schedule := scheduler.CorporationScheduler(corp)
		for _, s := range schedule {
			c.Schedule(s, cron.FuncJob(func() {
				if err := matcherService.PublishOffers(client, corp); err != nil {
					logger.Sugar().Error(err)
				}
			}))
		}
	}
}
