package job

import (
	"github.com/robfig/cron/v3"
	notificationsService "github.com/woningfinder/woningfinder/internal/services/notifications"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// SendWeeklyUpdate populates the weekly updates cron jobs (once a week)
func SendWeeklyUpdate(logger *logging.Logger, c *cron.Cron, userService userService.Service, notificationsService notificationsService.Service) {
	// emails update are send every fridays at 16:00
	spec := "0 0 16 * * 5"

	// populate cron
	c.AddJob(spec, cron.FuncJob(func() {
		// get all users
		users, err := userService.GetWeeklyUpdateUsers()
		if err != nil {
			logger.Sugar().Errorf("error while sending weekly update: %w", err)
		}

		// send confirmation email to each user
		for _, user := range users {
			if err := notificationsService.SendWeeklyUpdate(user, user.HousingPreferencesMatch); err != nil {
				logger.Sugar().Errorf("error while sending weekly update: %w", err)
			}
		}
	}))
}
