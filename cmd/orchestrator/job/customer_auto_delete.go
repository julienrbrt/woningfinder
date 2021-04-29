package job

import (
	"errors"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/robfig/cron/v3"
	"github.com/woningfinder/woningfinder/internal/customer"
	"github.com/woningfinder/woningfinder/internal/database"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// CustomerAutoDelete deletes the customers that did not paid within 24 hours
func CustomerAutoDelete(logger *logging.Logger, c *cron.Cron, userService userService.Service, dbClient database.DBClient) {
	// customer auto deletes are performed everyday at 00:00
	spec := "0 0 0 * * *"

	// populate cron
	c.AddJob(spec, cron.FuncJob(func() {
		var users []customer.User
		usersPlanQuery := dbClient.Conn().Model((*customer.UserPlan)(nil)).ColumnExpr("user_id")
		err := dbClient.Conn().
			Model(&users).
			Where("id NOT IN (?)", usersPlanQuery).
			Select()
		if err != nil && !errors.Is(err, pg.ErrNoRows) {
			logger.Sugar().Warnf("failed getting users to delete: %w", err)
		}

		for _, user := range users {
			// delete only user that did not paid since a day ago
			if user.CreatedAt.Before(time.Now().Add(-48 * time.Hour)) {
				if err := userService.DeleteUser(&user); err != nil {
					logger.Sugar().Warnf("failed deleting user %s: %w", user.Email, err)
				}
			}
		}
	}))
}
