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
	maxUnconfirmedTime      = 48 * time.Hour
)

// CustomerUnconfirmedCleanup reminds a unconfirmed customers to confirm their emails
// and deletes the customers that have an unconfirmed email for more than 48h
func (j *Jobs) CustomerUnconfirmedCleanup(c *cron.Cron) {
	// checks perfomed at 08:00, 16:00
	spec := "0 0 8,16 * * *"

	// populate cron
	c.AddJob(spec, cron.FuncJob(func() {
		j.logger.Sugar().Info("customer-unconfirmed-cleanup job started")

		var users []customer.User
		usersPlanQuery := j.dbClient.Conn().Model((*customer.UserPlan)(nil)).ColumnExpr("user_id")
		err := j.dbClient.Conn().
			Model(&users).
			Where("id NOT IN (?)", usersPlanQuery).
			Select()
		if err != nil && !errors.Is(err, pg.ErrNoRows) {
			j.logger.Sugar().Errorf("failed getting users get unconfirmed users: %w", err)
		}

		for _, user := range users {
			// check only user than did not paid since 4 hours
			switch {
			case user.CreatedAt.Before(time.Now().Add(-unconfirmedReminderTime)):
				// check if reminder already sent
				uuid := buildPaymentReminderUUID(&user)
				if j.isAlreadySent(uuid) {
					continue
				}

				// send reminder
				if err := j.emailService.SendEmailConfirmationReminder(&user); err != nil {
					j.logger.Sugar().Error("Error sending %s payment reminder: %w", user.Email, err)
				}

				// set reminder as sent
				j.markAsSent(uuid)

				// delete only user that did not paid since 48h
			case user.CreatedAt.Before(time.Now().Add(-maxUnconfirmedTime)):
				user, err := j.userService.GetUser(&user)
				if err != nil {
					j.logger.Sugar().Errorf("failed getting user %s for deleting: %w", user.Email, err)

				}

				if err := j.userService.DeleteUser(user); err != nil {
					j.logger.Sugar().Errorf("failed deleting user %s: %w", user.Email, err)
				}
			}
		}
	}))
}

func buildPaymentReminderUUID(user *customer.User) string {
	return base64.StdEncoding.EncodeToString([]byte(user.Email + "reminder email sent"))
}
