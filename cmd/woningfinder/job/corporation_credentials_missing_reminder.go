package job

import (
	"time"

	"github.com/julienrbrt/woningfinder/internal/customer"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var corporationCredentialsMissingReminderTime = map[int]time.Duration{
	1: 24 * time.Hour,  // 1 day
	2: 96 * time.Hour,  // 4 days
	3: 216 * time.Hour, // 9 days
	4: 336 * time.Hour, // 14 days
}

func (j *Jobs) SendCorporationCredentialsMissingReminder(c *cron.Cron) {
	// checks performed at 08:00, 16:00
	spec := "0 0 8,16 * * *"

	// populate cron
	c.AddJob(spec, cron.FuncJob(func() {
		j.logger.Info("send-corporation-credentials-missing-reminder job started")

		users, err := j.userService.GetWeeklyUpdateUsers()
		if err != nil {
			j.logger.Error("error while sending credentials missing reminder", zap.Error(err))
		}

		// send confirmation email to each user
		for _, user := range users {
			// skip inactive users and users with corporation credentials
			if !user.IsActivated() || len(user.CorporationCredentials) > 0 {
				continue
			}

			j.sendEmailCorporationCredentialsMissingReminder(user)
		}
	}))
}

// send reminder to users that still does not have corporation credentials
func (j *Jobs) sendEmailCorporationCredentialsMissingReminder(user *customer.User) {
	userReminderCount, err := j.userService.GetReminderCount(user.Email)
	if err != nil {
		j.logger.Error("error getting reminder counts", zap.String("email", user.Email), zap.Error(err))
		return
	}

	for count, duration := range corporationCredentialsMissingReminderTime {
		// check if reminder already sent
		if userReminderCount.CorporationCredentialsMissingReminderCount >= count {
			continue
		}

		if time.Until(user.CreatedAt.Add(duration)) <= 0 {
			// send reminder
			if err := j.emailService.SendCorporationCredentialsMissing(user); err != nil {
				j.logger.Error("error while sending credentials missing reminder", zap.Error(err))
				return
			}

			// update reminder count
			userReminderCount.CorporationCredentialsMissingReminderCount = count
			if err := j.userService.UpdateReminderCount(user.Email, userReminderCount); err != nil {
				j.logger.Error("error updating reminder count", zap.String("email", user.Email), zap.Error(err))
			}

			return
		}
	}
}
