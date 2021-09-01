package job

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/robfig/cron/v3"
	"github.com/woningfinder/woningfinder/internal/customer"
)

const (
	unconfirmedReminderTime = 4 * time.Hour
	maxUnconfirmedTime      = 72 * time.Hour
)

// CustomerUnconfirmedCleanup reminds a unconfirmed customers to confirm their emails
// and deletes the customers that have an unconfirmed email for more than unconfirmedReminderTime
func (j *Jobs) CustomerUnconfirmedCleanup(c *cron.Cron) {
	// checks perfomed at 08:00, 16:00
	spec := "0 0 8,16 * * *"

	// populate cron
	c.AddJob(spec, cron.FuncJob(func() {
		j.logger.Sugar().Info("customer-unconfirmed-cleanup job started")

		var users []*customer.User
		// delete unconfirmed account
		err := j.dbClient.Conn().
			Model(&users).
			Join("INNER JOIN user_plans up ON id = up.user_id").
			Where("free_trial_started_at IS NULL").
			Select()
		if err != nil && !errors.Is(err, pg.ErrNoRows) {
			j.logger.Sugar().Errorf("failed getting users get unconfirmed users: %w", err)
		}

		for _, user := range users {
			// send reminder to users that didn't confirm their email
			switch {
			case time.Until(user.CreatedAt.Add(unconfirmedReminderTime)) <= 0:
				// check if reminder already sent
				uuid := buildCustomerUnconfirmedEmailReminderUUID(user)
				if j.isAlreadySent(uuid) {
					continue
				}

				// send reminder
				if err := j.emailService.SendEmailConfirmationReminder(user); err != nil {
					j.logger.Sugar().Error("error sending %s email confirmation reminder: %w", user.Email, err)
				}

				// set reminder as sent
				j.markAsSent(uuid)

			// delete only user that did not paid since 48h
			case time.Until(user.CreatedAt.Add(maxUnconfirmedTime)) <= 0:
				if err := j.userService.DeleteUser(user.Email); err != nil {
					j.logger.Sugar().Errorf("failed deleting user %s: %w", user.Email, err)
				}
			}
		}
	}))
}

func buildCustomerUnconfirmedEmailReminderUUID(user *customer.User) string {
	return base64.StdEncoding.EncodeToString([]byte(user.Email + "customer confirmation email reminder sent"))
}
