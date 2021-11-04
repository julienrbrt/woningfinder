package job

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func (j *Jobs) ReminderEmail(c *cron.Cron) {
	// checks performed at 08:00, 16:00
	spec := "0 0 8,16 * * *"

	// populate cron
	c.AddJob(spec, cron.FuncJob(func() {
		// reminder no corporation credentials

	}))
}

func (j *Jobs) CorporationCredentialsMissingReminder() {
	users, err := j.userService.GetWeeklyUpdateUsers()
	if err != nil {
		j.logger.Error("error while sending weekly update", zap.Error(err))
	}

	// send confirmation email to each user
	for _, user := range users {
		// skip user with invalid plan
		if !user.Plan.IsValid() {
			continue
		}

	}
}
