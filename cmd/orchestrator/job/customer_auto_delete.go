package job

import (
	"github.com/robfig/cron/v3"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// CustomerAutoDelete deletes the customers that did not paid within 24 hours
func CustomerAutoDelete(logger *logging.Logger, c *cron.Cron, userService userService.Service) {
	// customer auto deletes are performed everyday at 00:00
	spec := "0 0 0 * * *"

	// populate cron
	c.AddJob(spec, cron.FuncJob(func() {
		logger.Sugar().Info("customer auto delete ran but not implemented.")
	}))
}
