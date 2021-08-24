package user

import (
	"errors"
	"fmt"
	"time"

	pg "github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/woningfinder/woningfinder/internal/customer"
)

// TODO eventually use a prepare function to create it in one query only
func (s *service) CreateUser(u *customer.User) error {
	db := s.dbClient.Conn()

	// verify user
	if err := u.HasMinimal(); err != nil {
		return fmt.Errorf("error user %s invalid: %w", u.Email, err)
	}

	// a user cannot have paid when being created so reset by security
	u.Plan.CreatedAt = (time.Time{})
	u.Plan.PurchasedAt = (time.Time{})

	// create user - if exist throw error
	if _, err := db.Model(u).Insert(); err != nil {
		return fmt.Errorf("failing creating user: %w", err)
	}

	// create user housing preferences
	if err := s.CreateHousingPreferences(u, u.HousingPreferences); err != nil {
		// rollback
		if _, err2 := db.Model(u).Where("email ILIKE ?", u.Email).Delete(); err2 != nil {
			s.logger.Sugar().Errorf("error %w and error when rolling back user creation: %w", err, err2)
		}

		return fmt.Errorf("error when creating user %s: %w", u.Email, err)
	}

	return nil
}

func (s *service) GetUser(search *customer.User) (*customer.User, error) {
	db := s.dbClient.Conn()

	var user customer.User
	// get user (by id or email)
	if search.ID > 0 {
		if err := db.Model(&user).Where("id = ?", search.ID).Select(); err != nil {
			return nil, fmt.Errorf("failed getting user %d: %w", search.ID, err)
		}
	} else {
		if err := db.Model(&user).Where("email ILIKE ?", search.Email).Select(); err != nil {
			return nil, fmt.Errorf("failed getting user %s: %w", search.Email, err)
		}
	}

	// enrich user
	var err error
	user.HousingPreferences, err = s.GetHousingPreferences(&user)
	if err != nil {
		return nil, fmt.Errorf("failed getting user %s housing preferences: %w", user.Email, err)
	}

	if err := db.Model(&user).Where("id = ?", user.ID).Relation("Plan").Select(); err != nil {
		return nil, fmt.Errorf("failed getting user %s plan: %w", user.Email, err)
	}

	if err := db.Model(&user).Where("id = ?", user.ID).Relation("HousingPreferencesMatch").Select(); err != nil {
		return nil, fmt.Errorf("failed getting user %s housing preferences match: %w", user.Email, err)
	}

	// enriching with corporation credentials is skipped because not useful
	return &user, nil
}

// ConfirmUser validate user account and start free trial
func (s *service) ConfirmUser(email string, plan customer.Plan) (*customer.User, error) {
	// get user
	user, err := s.GetUser(&customer.User{Email: email})
	if err != nil {
		return nil, fmt.Errorf("error while processing payment data: cannot get user: %w", err)
	}

	// set that has started its free trial
	if _, err := s.dbClient.Conn().
		Model(&customer.UserPlan{UserID: user.ID, PlanName: plan.Name}).
		OnConflict("(user_id) DO UPDATE").
		Insert(); err != nil {
		return nil, fmt.Errorf("error when adding user plan: %w", err)
	}

	return user, nil
}

// ConfirmPayment set that the user has paid
func (s *service) ConfirmPayment(email string, plan customer.Plan) (*customer.User, error) {
	// get user
	user, err := s.GetUser(&customer.User{Email: email})
	if err != nil {
		return nil, fmt.Errorf("error while processing payment data: cannot get user: %w", err)
	}

	// set that user has paid
	if _, err := s.dbClient.Conn().
		Model(&customer.UserPlan{UserID: user.ID, PlanName: plan.Name, PurchasedAt: time.Now()}).
		OnConflict("(user_id) DO UPDATE").
		Insert(); err != nil {
		return nil, fmt.Errorf("error when adding user plan: %w", err)
	}

	return user, nil
}

// TODO eventually use a prepare function to create it in one query only and improve performance
func (s *service) GetWeeklyUpdateUsers() ([]*customer.User, error) {
	db := s.dbClient.Conn()

	userList := []customer.UserPlan{}
	if err := db.Model(&userList).Order("created_at ASC").Select(); err != nil {
		return nil, fmt.Errorf("error getting paid users list: %w", err)
	}

	var users []*customer.User
	for _, user := range userList {
		u := &customer.User{ID: user.UserID}

		if err := db.Model(u).Where("id = ?", u.ID).Select(); err != nil {
			return nil, fmt.Errorf("failed getting user %d: %w", u.ID, err)
		}

		// enrich housing preferences match
		if err := db.Model(u).
			Where("id = ?", u.ID).
			Relation("HousingPreferencesMatch", func(q *orm.Query) (*orm.Query, error) {
				return q.Where("created_at >= now() - interval '7 day'"), nil
			}).
			Select(); err != nil {
			return nil, fmt.Errorf("failed getting user %s housing preferences match: %w", u.Email, err)
		}

		// enrich corporation credentials
		if err := db.Model(&u.CorporationCredentials).Where("user_id = ?", u.ID).Select(); err != nil {
			return nil, fmt.Errorf("error when getting corporation credentials for userID %d: %w", u.ID, err)
		}

		users = append(users, u)
	}

	return users, nil
}

// TODO eventually use a prepare function to create it in one query only and improve performance
func (s *service) DeleteUser(u *customer.User) error {
	db := s.dbClient.Conn()

	// delete all corporations credentials
	if _, err := db.Model((*customer.CorporationCredentials)(nil)).Where("user_id = ?", u.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed deleting housing corporation credentials for %s: %w", u.Email, err)
	}

	// delete housing preferences
	if _, err := db.Model((*customer.HousingPreferencesCity)(nil)).Where("housing_preferences_id = ?", u.HousingPreferences.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed deleting housing preferences cities for %s: %w", u.Email, err)
	}

	if _, err := db.Model((*customer.HousingPreferencesCityDistrict)(nil)).Where("housing_preferences_id = ?", u.HousingPreferences.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed deleting housing preferences cities districts for %s: %w", u.Email, err)
	}

	if _, err := db.Model((*customer.HousingPreferencesHousingType)(nil)).Where("housing_preferences_id = ?", u.HousingPreferences.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed deleting housing preferences housing type for %s: %w", u.Email, err)
	}

	if _, err := db.Model((*customer.HousingPreferences)(nil)).Where("user_id = ?", u.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed deleting housing preferences cities for %s: %w", u.Email, err)
	}

	if _, err := db.Model((*customer.HousingPreferencesMatch)(nil)).Where("user_id = ?", u.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed deleting housing preferences match for %s: %w", u.Email, err)
	}

	// delete user
	if _, err := db.Model((*customer.UserPlan)(nil)).Where("user_id = ?", u.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed delete user plan for %s: %w", u.Email, err)
	}

	if _, err := db.Model((*customer.User)(nil)).Where("id = ?", u.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed delete user for %s: %w", u.Email, err)
	}

	return nil
}
