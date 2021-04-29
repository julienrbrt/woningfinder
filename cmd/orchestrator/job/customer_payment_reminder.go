package job

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/robfig/cron/v3"
	"github.com/woningfinder/woningfinder/internal/customer"
	"github.com/woningfinder/woningfinder/internal/database"
)

// SendCustomerPaymentReminder reminds a unpaid customer to complete it's housing preferences (aka payment)
func (j *Jobs) SendCustomerPaymentReminder(c *cron.Cron) {
	// customer auto deletes are performed twice a day at 12:00 and 21:00
	spec := "0 0 12,21 * * *"

	// populate cron
	c.AddJob(spec, cron.FuncJob(func() {
		j.logger.Sugar().Info("send-customer-payment-reminder job started")

		var users []customer.User
		usersPlanQuery := j.dbClient.Conn().Model((*customer.UserPlan)(nil)).ColumnExpr("user_id")
		err := j.dbClient.Conn().
			Model(&users).
			Where("id NOT IN (?)", usersPlanQuery).
			Select()
		if err != nil && !errors.Is(err, pg.ErrNoRows) {
			j.logger.Sugar().Warnf("failed getting users to send payment reminder: %w", err)
		}

		for _, user := range users {
			// check only user than did not paid since 8 hours
			if user.CreatedAt.Before(time.Now().Add(-8 * time.Hour)) {
				// check if reminder already sent
				uuid := buildPaymentReminderUUID(&user)
				if j.hasPaymentReminder(uuid) {
					continue
				}

				// send reminder
				if err := j.notificationsService.SendPaymentReminder(&user); err != nil {
					j.logger.Sugar().Error("Error sending %s payment reminder: %w", user.Email, err)
				}

				// set reminder as sent
				j.setPaymentReminderSent(uuid)
			}
		}
	}))
}

// hasPaymentReminder check if a user already received a  payment reminder
func (j *Jobs) hasPaymentReminder(uuid string) bool {
	_, err := j.redisClient.Get(uuid)
	if err != nil {
		if !errors.Is(err, database.ErrRedisKeyNotFound) {
			j.logger.Sugar().Errorf("error when getting reaction: %w", err)
		}
		// does not have received reminder
		return false
	}

	return true
}

// setPaymentReminderSent saves that an user has received a payment reminder
func (j *Jobs) setPaymentReminderSent(uuid string) {
	if err := j.redisClient.Set(uuid, true); err != nil {
		j.logger.Sugar().Errorf("error when saving reaction to redis: %w", err)
	}
}

func buildPaymentReminderUUID(user *customer.User) string {
	return base64.StdEncoding.EncodeToString([]byte(user.Email + "reminder email sent"))
}
