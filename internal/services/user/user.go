package user

import (
	"fmt"
	"time"

	"github.com/go-pg/pg/v10/orm"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

// TODO eventually use a prepare function to create it in one query only
func (s *service) CreateUser(u *entity.User) error {
	db := s.dbClient.Conn()

	// verify user
	if err := u.IsValid(); err != nil {
		return fmt.Errorf("error user %s invalid: %w", u.Email, err)
	}

	// a user cannot have paid when being created so reset by security
	u.Plan.CreatedAt = (time.Time{})

	// create user - if exist throw error
	if _, err := db.Model(u).Insert(); err != nil {
		return fmt.Errorf("failing creating user: %w", err)
	}

	// create user housing preferences
	if err := s.CreateHousingPreferences(u, u.HousingPreferences); err != nil {
		// rollback
		if _, err2 := db.Model(u).Where("email = ?", u.Email).Delete(); err2 != nil {
			s.logger.Sugar().Errorf("error %w and error when rolling back user creation: %w", err, err2)
		}

		return fmt.Errorf("error when creating user %s: %w", u.Email, err)
	}

	return nil
}

func (s *service) GetUser(search *entity.User) (*entity.User, error) {
	db := s.dbClient.Conn()

	var u entity.User
	// get user (by id or email)
	if search.ID > 0 {
		if err := db.Model(&u).Where("id = ?", search.ID).Select(); err != nil {
			return nil, fmt.Errorf("failed getting user %d: %w", search.ID, err)
		}
	} else {
		if err := db.Model(&u).Where("email ILIKE ?", search.Email).Select(); err != nil {
			return nil, fmt.Errorf("failed getting user %s: %w", search.Email, err)
		}
	}

	// enrich user
	var err error
	u.HousingPreferences, err = s.GetHousingPreferences(&u)
	if err != nil {
		return nil, fmt.Errorf("failed getting user %s housing preferences: %w", u.Email, err)
	}

	if err := db.Model(&u).Where("id = ?", u.ID).Relation("Plan").Select(); err != nil {
		return nil, fmt.Errorf("failed getting user %s plan: %w", u.Email, err)
	}

	if err := db.Model(&u).Where("id = ?", u.ID).Relation("HousingPreferencesMatch").Select(); err != nil {
		return nil, fmt.Errorf("failed getting user %s housing preferences match: %w", u.Email, err)
	}

	// enriching with corporation credentials is skipped because not useful

	return &u, nil
}

// TODO eventually use a prepare function to create it in one query only and improve performance
func (s *service) GetWeeklyUpdateUsers() ([]*entity.User, error) {
	db := s.dbClient.Conn()

	userList := []entity.UserPlan{}
	if err := db.Model(&userList).Order("created_at ASC").Select(); err != nil {
		return nil, fmt.Errorf("error getting paid users list: %w", err)
	}

	var users []*entity.User
	for _, user := range userList {
		u := &entity.User{ID: user.UserID}

		if err := db.Model(u).Where("id = ?", u.ID).Select(); err != nil {
			return nil, fmt.Errorf("failed getting user %d: %w", u.ID, err)
		}

		if err := db.Model(u).
			Where("id = ?", u.ID).
			Relation("HousingPreferencesMatch", func(q *orm.Query) (*orm.Query, error) {
				return q.Where("created_at >= now() - interval '1 day'"), nil
			}).
			Select(); err != nil {
			return nil, fmt.Errorf("failed getting user %s housing preferences match: %w", u.Email, err)
		}

		users = append(users, u)
	}

	return users, nil
}

func (s *service) DeleteUser(u *entity.User) error {
	// TODO to implement
	// delete all corporations credentials
	// delete housing preferences
	// delete user
	panic("not implemented")
}
