package job

import (
	"errors"
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
	userReminderCount, err := j.userService.GetReminderCount(user.Email)
	if err != nil {
		j.logger.Error("error getting reminder counts", zap.String("email", user.Email), zap.Error(err))
		return
	}

	for count, duration := range unconfirmedReminderTime {
		// check if reminder already sent
		if userReminderCount.UnconfirmedReminderCount >= count {
			continue
		}

		if time.Until(user.CreatedAt.Add(duration)) <= 0 {
			// send reminder
			if err := j.emailService.SendEmailConfirmationReminder(user); err != nil {
				j.logger.Error("error sending email confirmation reminder", zap.String("email", user.Email), zap.Error(err))
				return
			}

			// update reminder count
			userReminderCount.UnconfirmedReminderCount = count
			if err := j.userService.UpdateReminderCount(user.Email, userReminderCount); err != nil {
				j.logger.Error("error updating reminder count", zap.String("email", user.Email), zap.Error(err))
			}

			return
		}
	}
}
