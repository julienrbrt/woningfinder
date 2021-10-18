package job

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/robfig/cron/v3"
	"github.com/woningfinder/woningfinder/internal/customer"
)

const (
	firstUnconfirmedReminderTime  = 24 * time.Hour
	secondUnconfirmedReminderTime = 72 * time.Hour
	maxUnconfirmedTime            = 30 * 24 * time.Hour
)

// CleanupUnconfirmedCustomer reminds a unconfirmed customers to confirm their account
// and deletes the customers that have an unconfirmed email for more than maxUnconfirmedTime
func (j *Jobs) CleanupUnconfirmedCustomer(c *cron.Cron) {
	// checks perfomed at 08:00, 16:00
	spec := "0 0 8,16 * * *"

	// populate cron
	c.AddJob(spec, cron.FuncJob(func() {
		j.logger.Sugar().Info("cleanup-unconfirmed-customer job started")

		var users []*customer.User
		// delete unconfirmed account
		err := j.dbClient.Conn().
			Model(&users).
			Join("INNER JOIN user_plans up ON id = up.user_id").
			Where("up.activated_at IS NULL").
			Select()
		if err != nil && !errors.Is(err, pg.ErrNoRows) {
			j.logger.Sugar().Errorf("failed getting unconfirmed users: %w", err)
		}

		for _, user := range users {
			// send reminder to users that didn't confirm their email
			if time.Until(user.CreatedAt.Add(firstUnconfirmedReminderTime)) <= 0 {
				j.sendEmailReminder(user, 1)
			}

			if time.Until(user.CreatedAt.Add(secondUnconfirmedReminderTime)) <= 0 {
				j.sendEmailReminder(user, 2)
			}

			// delete only user that did confirm their email since 30 days
			if time.Until(user.CreatedAt.Add(maxUnconfirmedTime)) <= 0 {
				if err := j.userService.DeleteUser(user.Email); err != nil {
					j.logger.Sugar().Errorf("failed deleting user %s: %w", user.Email, err)
				}
			}
		}
	}))
}

func (j *Jobs) sendEmailReminder(user *customer.User, count int) {
	// check if reminder already sent
	uuid := base64.StdEncoding.EncodeToString([]byte(user.Email + fmt.Sprintf("customer confirmation email reminder %d sent", count)))
	if j.redisClient.HasUUID(uuid) {
		return
	}

	// send reminder
	if err := j.emailService.SendEmailConfirmationReminder(user); err != nil {
		j.logger.Sugar().Error("error sending %s email confirmation reminder: %w", user.Email, err)
	}

	// set reminder as sent
	j.redisClient.SetUUID(uuid)
}