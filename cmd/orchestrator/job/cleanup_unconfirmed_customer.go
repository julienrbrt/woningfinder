package job

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/julienrbrt/woningfinder/internal/customer"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

const maxUnconfirmedTime = 30 * 24 * time.Hour

var unconfirmedReminderTime = map[int]time.Duration{
	1: 24 * time.Hour, // 1 day
	2: 72 * time.Hour, // 3 days
}

// CleanupUnconfirmedCustomer reminds a unconfirmed customers to confirm their account
// and deletes the customers that have an unconfirmed email for more than maxUnconfirmedTime
func (j *Jobs) CleanupUnconfirmedCustomer(c *cron.Cron) {
	// checks perfomed at 08:00, 16:00
	spec := "0 0 8,16 * * *"

	// populate cron
	c.AddJob(spec, cron.FuncJob(func() {
		j.logger.Info("cleanup-unconfirmed-customer job started")

		var users []*customer.User
		// delete unconfirmed account
		err := j.dbClient.Conn().
			Model(&users).
			Where("activated_at IS NULL").
			Select()
		if err != nil && !errors.Is(err, pg.ErrNoRows) {
			j.logger.Error("failed getting unconfirmed users", zap.Error(err))
		}

		for _, user := range users {
			j.sendEmailConfirmationReminder(user)

			// delete only unsubcribed user that did confirm their email since 30 days
			if !user.IsActivated() && time.Until(user.CreatedAt.Add(maxUnconfirmedTime)) <= 0 {
				if err := j.userService.DeleteUser(user.Email); err != nil {
					j.logger.Error("failed deleting user", zap.String("email", user.Email), zap.Error(err))
				}
			}
		}
	}))
}

// send reminder to users that didn't confirm their email
func (j *Jobs) sendEmailConfirmationReminder(user *customer.User) {
	for count, duration := range unconfirmedReminderTime {
		// check if reminder already sent
		uuid := hex.EncodeToString([]byte(user.Email + fmt.Sprintf("customer confirmation email reminder %d sent", count)))
		if j.redisClient.HasUUID(uuid) {
			continue
		}

		if time.Until(user.CreatedAt.Add(duration)) <= 0 {
			// send reminder
			if err := j.emailService.SendEmailConfirmationReminder(user); err != nil {
				j.logger.Error("error sending email confirmation reminder", zap.String("email", user.Email), zap.Error(err))
				return
			}

			j.redisClient.SetUUID(uuid)
			return
		}
	}
}
