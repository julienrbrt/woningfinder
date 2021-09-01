package job

import (
	"github.com/robfig/cron/v3"
)

// SendWeeklyUpdate populates the weekly updates cron jobs (once a week)
func (j *Jobs) SendWeeklyUpdate(c *cron.Cron) {
	// emails update are send every friday at 20:00
	spec := "0 0 20 * * 5"

	// populate cron
	c.AddJob(spec, cron.FuncJob(func() {
		j.logger.Sugar().Info("send-weekly-update job started")

		users, err := j.userService.GetUsersWithHousingPreferencesMatch()
		if err != nil {
			j.logger.Sugar().Errorf("error while sending weekly update: %w", err)
		}

		// send confirmation email to each user
		for _, user := range users {
			// user has no corporation credentials and no match didn't react for them be we cannot
			if len(user.HousingPreferencesMatch) == 0 && len(user.CorporationCredentials) == 0 {
				if err := j.emailService.SendCorporationCredentialsMissing(user); err != nil {
					j.logger.Sugar().Errorf("error while sending weekly update (credentials missing): %w", err)
				}
				continue
			}

			if err := j.emailService.SendWeeklyUpdate(user, user.HousingPreferencesMatch); err != nil {
				j.logger.Sugar().Errorf("error while sending weekly update: %w", err)
			}
		}
	}))
}
