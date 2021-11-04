package job

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/robfig/cron/v3"
	"github.com/woningfinder/woningfinder/internal/customer"
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
			Relation("Plan").
			Join("INNER JOIN user_plans up ON id = up.user_id").
			Where("up.activated_at IS NULL").
			Select()
		if err != nil && !errors.Is(err, pg.ErrNoRows) {
			j.logger.Error("failed getting unconfirmed users", zap.Error(err))
		}

		for _, user := range users {
			// send reminder to users that didn't confirm their email
			for count, duration := range unconfirmedReminderTime {
				if time.Until(user.CreatedAt.Add(duration)) <= 0 {
					j.sendEmailConfirmationReminder(user, count)
				}
			}

			// delete only unsubcribed user that did confirm their email since 30 days
			if !user.Plan.IsSubscribed() && time.Until(user.CreatedAt.Add(maxUnconfirmedTime)) <= 0 {
				if err := j.userService.DeleteUser(user.Email); err != nil {
					j.logger.Error("failed deleting user", zap.String("email", user.Email), zap.Error(err))
				}
			}
		}
	}))
}

func (j *Jobs) sendEmailConfirmationReminder(user *customer.User, count int) {
	// check if reminder already sent
	uuid := base64.StdEncoding.EncodeToString([]byte(user.Email + fmt.Sprintf("customer confirmation email reminder %d sent", count)))
	if j.redisClient.HasUUID(uuid) {
		return
	}

	// send reminder
	if err := j.emailService.SendEmailConfirmationReminder(user); err != nil {
		j.logger.Error("error sending email confirmation reminder", zap.String("email", user.Email), zap.Error(err))
		return
	}

	// set reminder as sent
	j.redisClient.SetUUID(uuid)
}
