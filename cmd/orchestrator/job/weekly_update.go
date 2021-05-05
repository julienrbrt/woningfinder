package job

import (
	"github.com/robfig/cron/v3"
)

// SendWeeklyUpdate populates the weekly updates cron jobs (once a week)
func (j *Jobs) SendWeeklyUpdate(c *cron.Cron) {
	// emails update are send every fridays at 16:00
	spec := "0 0 16 * * 5"

	// populate cron
	c.AddJob(spec, cron.FuncJob(func() {
		j.logger.Sugar().Info("send-weekly-update job started")

		// get all users
		users, err := j.userService.GetWeeklyUpdateUsers()
		if err != nil {
			j.logger.Sugar().Errorf("error while sending weekly update: %w", err)
		}

		// send confirmation email to each user
		for _, user := range users {
			if err := j.notificationService.SendWeeklyUpdate(user, user.HousingPreferencesMatch); err != nil {
				j.logger.Sugar().Errorf("error while sending weekly update: %w", err)
			}
		}
	}))
}
