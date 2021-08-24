package job

import (
	"encoding/base64"
	"errors"

	"github.com/go-pg/pg/v10"
	"github.com/robfig/cron/v3"
	"github.com/woningfinder/woningfinder/internal/customer"
)

// SendCustomerEndFreeTrialReminder reminds a free trial customer to pay the plan
func (j *Jobs) SendCustomerEndFreeTrialReminder(c *cron.Cron) {
	// checks performed everyday at 08:00, 14:00, 20:00
	spec := "0 0 8,14,20 * * *"

	// populate cron
	c.AddJob(spec, cron.FuncJob(func() {
		j.logger.Sugar().Info("send-customer-end-free-trial-reminder job started")

		// get all users without
		freeTrialPlans := j.dbClient.Conn().
			Model((*customer.UserPlan)(nil)).
			Where("purchased_at IS NULL").
			ColumnExpr("user_id").
			Select()

		var users []customer.User
		err := j.dbClient.Conn().
			Model(&users).
			Where("id IN (?)", freeTrialPlans).
			Select()
		if err != nil && !errors.Is(err, pg.ErrNoRows) {
			j.logger.Sugar().Errorf("failed getting users to send free trial reminder: %w", err)
		}

		for _, user := range users {
			if user.Plan.IsValid() {
				continue
			}

			// check if reminder already sent
			uuid := buildFreeTrialReminderReminderUUID(&user)
			if j.isAlreadySent(uuid) {
				continue
			}

			// send reminder
			if err := j.emailService.SendFreeTrialReminder(&user); err != nil {
				j.logger.Sugar().Error("Error sending %s free trial reminder: %w", user.Email, err)
			}

			// set reminder as sent
			j.markAsSent(uuid)
		}
	}))
}

func buildFreeTrialReminderReminderUUID(user *customer.User) string {
	return base64.StdEncoding.EncodeToString([]byte(user.Email + "free trial email reminder sent"))
}
