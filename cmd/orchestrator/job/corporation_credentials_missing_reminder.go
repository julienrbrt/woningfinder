package job

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/woningfinder/woningfinder/internal/customer"
	"go.uber.org/zap"
)

var corporationCredentialsMissingReminderTime = map[int]time.Duration{
	1: 24 * time.Hour,  // 1 day
	2: 72 * time.Hour,  // 3 days
	3: 168 * time.Hour, // 7 days
	4: 240 * time.Hour, // 10 days
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
			// skip user with invalid plan and users with corporation credentials
			if !user.Plan.IsValid() || len(user.CorporationCredentials) > 0 {
				continue
			}

			// send reminder to users that still does not have corporation credentials
			for count, duration := range corporationCredentialsMissingReminderTime {
				if time.Until(user.CreatedAt.Add(duration)) <= 0 {
					j.sendEmailCorporationCredentialsMissingReminder(user, count)
				}
			}
		}
	}))
}

func (j *Jobs) sendEmailCorporationCredentialsMissingReminder(user *customer.User, count int) {
	// check if reminder already sent
	uuid := base64.StdEncoding.EncodeToString([]byte(user.Email + fmt.Sprintf("customer corporation credentials missing reminder %d sent", count)))
	if j.redisClient.HasUUID(uuid) {
		return
	}

	// send reminder
	if err := j.emailService.SendCorporationCredentialsMissing(user); err != nil {
		j.logger.Error("error while sending credentials missing reminder", zap.Error(err))
		return
	}

	// set reminder as sent
	j.redisClient.SetUUID(uuid)
}
