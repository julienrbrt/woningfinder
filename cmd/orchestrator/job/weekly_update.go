package job

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// SendWeeklyUpdate populates the weekly updates cron jobs (once a week)
func (j *Jobs) SendWeeklyUpdate(c *cron.Cron) {
	// emails update are send every friday at 18:00
	spec := "0 0 18 * * 5"

	// populate cron
	c.AddJob(spec, cron.FuncJob(func() {
		j.logger.Info("send-weekly-update job started")

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

			// user has no corporation credentials and no match didn't react for them be we cannot send weekly update
			if len(user.HousingPreferencesMatch) == 0 && len(user.CorporationCredentials) == 0 {
				if err := j.emailService.SendCorporationCredentialsMissing(user); err != nil {
					j.logger.Error("error while sending weekly update (credentials missing)", zap.Error(err))
				}
				continue
			}

			if err := j.emailService.SendWeeklyUpdate(user, user.HousingPreferencesMatch); err != nil {
				j.logger.Error("error while sending weekly update", zap.Error(err))
			}
		}
	}))
}
