package job

import (
	"encoding/base64"
	"errors"

	"github.com/go-pg/pg/v10"
	"github.com/robfig/cron/v3"
	"github.com/woningfinder/woningfinder/internal/customer"
)

// SendCustomerEndFreeTrialReminder reminds a free trial customer to pay its plan
func (j *Jobs) SendCustomerEndFreeTrialReminder(c *cron.Cron) {
	// checks performed everyday at 08:00, 14:00, 20:00
	spec := "0 0 8,14,20 * * *"

	// populate cron
	c.AddJob(spec, cron.FuncJob(func() {
		j.logger.Sugar().Info("send-customer-end-free-trial-reminder job started")

		var users []*customer.User
		// get all users without paid account
		err := j.dbClient.Conn().
			Model(&users).
			Relation("Plan").
			Join("INNER JOIN user_plans up ON id = up.user_id").
			Where("up.free_trial_started_at IS NOT NULL AND up.purchased_at IS NULL").
			Select()
		if err != nil && !errors.Is(err, pg.ErrNoRows) {
			j.logger.Sugar().Errorf("failed getting users to check free trial: %w", err)
		}

		for _, user := range users {
			// skip paid users and users with valid free trial
			if user.Plan.IsValid() {
				continue
			}

			// check if reminder already sent
			uuid := buildFreeTrialReminderReminderUUID(user)
			if !j.redisClient.HasUUID(uuid) {

				// send reminder
				if err := j.emailService.SendFreeTrialReminder(user); err != nil {
					j.logger.Sugar().Error("error sending %s free trial reminder: %w", user.Email, err)
				}

				// set reminder as sent
				j.redisClient.SetUUID(uuid)
			}
		}
	}))
}

func buildFreeTrialReminderReminderUUID(user *customer.User) string {
	return base64.StdEncoding.EncodeToString([]byte(user.Email + "free trial email reminder sent"))
}
