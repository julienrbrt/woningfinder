package job

import (
	"errors"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/robfig/cron/v3"
	"github.com/woningfinder/woningfinder/internal/customer"
)

// CustomerAutoDelete deletes the customers that did not paid within 24 hours
func (j *Jobs) CustomerAutoDelete(c *cron.Cron) {
	// customer auto deletes are performed everyday at 00:00
	spec := "0 0 0 * * *"

	// populate cron
	c.AddJob(spec, cron.FuncJob(func() {
		j.logger.Sugar().Info("customer-auto-delete job started")

		var users []customer.User
		usersPlanQuery := j.dbClient.Conn().Model((*customer.UserPlan)(nil)).ColumnExpr("user_id")
		err := j.dbClient.Conn().
			Model(&users).
			Where("id NOT IN (?)", usersPlanQuery).
			Select()
		if err != nil && !errors.Is(err, pg.ErrNoRows) {
			j.logger.Sugar().Warnf("failed getting users to delete: %w", err)
		}

		for _, user := range users {
			// delete only user that did not paid since a day ago
			if user.CreatedAt.Before(time.Now().Add(-48 * time.Hour)) {
				if err := j.userService.DeleteUser(&user); err != nil {
					j.logger.Sugar().Warnf("failed deleting user %s: %w", user.Email, err)
				}
			}
		}
	}))
}
