package user

import (
	"errors"
	"fmt"

	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/customer"
	"github.com/woningfinder/woningfinder/internal/services"
)

func (s *service) CreateCorporationCredentials(userID uint, credentials customer.CorporationCredentials) error {
	if credentials.CorporationName == "" {
		return errors.New("error when creating corporation credentials: corporation invalid")
	}

	if credentials.Login == "" || credentials.Password == "" {
		return errors.New("error when creating corporation credentials: login or password missing")
	}

	// check credentials validity
	if err := s.validateCredentials(credentials); err != nil {
		return fmt.Errorf("error when validation corporation credentials: %w", err)
	}

	// encrypt credentials
	var err error
	secretKey := auth.BuildAESKey(userID, credentials.CorporationName, s.aesSecret)
	credentials.Login, err = auth.EncryptString(credentials.Login, secretKey)
	if err != nil {
		return fmt.Errorf("error when encrypting corporation credentials: %w", err)
	}

	credentials.Password, err = auth.EncryptString(credentials.Password, secretKey)
	if err != nil {
		return fmt.Errorf("error when encrypting corporation credentials: %w", err)
	}

	// assign ids
	credentials.UserID = userID

	// check if already existing
	if _, err := s.dbClient.Conn().Model(&credentials).
		OnConflict("(user_id, corporation_name) DO UPDATE").
		Insert(); err != nil {
		return fmt.Errorf("error when creating or updating corporation credentials: %w", err)
	}

	return nil
}

func (s *service) GetCorporationCredentials(userID uint, corporationName string) (*customer.CorporationCredentials, error) {
	credentials := customer.CorporationCredentials{
		UserID:          userID,
		CorporationName: corporationName,
	}

	// get corporation credentials
	if err := s.dbClient.Conn().Model(&credentials).Where("user_id = ? and corporation_name = ?", credentials.UserID, corporationName).Select(); err != nil {
		return nil, fmt.Errorf("error when getting corporation credentials for userID %d: %w", userID, err)
	}

	return &credentials, nil
}

// GetAllCorporationCredentials gets all the user corporation credentials that we need to react for
// The corporation credentials are sorted by account creation date (and not by plan - see issue #21)
func (s *service) GetAllCorporationCredentials(corporationName string) ([]customer.CorporationCredentials, error) {
	credentials := []customer.CorporationCredentials{}
	if err := s.dbClient.Conn().
		Model(&credentials).
		Where("corporation_name = ?", corporationName).
		Order("created_at ASC"). // people having registered their credentials for the longer get reaction priority
		Select(); err != nil {
		return nil, fmt.Errorf("error getting user credentials: %w", err)
	}

	// no users found
	if len(credentials) == 0 {
		return nil, services.ErrNoMatchFound
	}

	return credentials, nil
}

func (s *service) UpdateCorporationCredentialsFailureCount(userID uint, corporationName string, failureCount int) error {
	credentials := customer.CorporationCredentials{
		UserID:          userID,
		CorporationName: corporationName,
	}

	// update failure count
	if _, err := s.dbClient.Conn().Model(&credentials).
		Set("failure_count = ?", failureCount).
		Where("user_id = ? and corporation_name = ?", credentials.UserID, corporationName).
		Update(); err != nil {
		return fmt.Errorf("error whenupdating corporation credentials failure count: %w", err)
	}

	return nil
}

func (s *service) DeleteCorporationCredentials(userID uint, corporationName string) error {
	credentials, err := s.GetCorporationCredentials(userID, corporationName)
	if err != nil {
		return fmt.Errorf("error when getting corporation credentials: %w", err)
	}

	// delete permanently
	if _, err := s.dbClient.Conn().Model((*customer.CorporationCredentials)(nil)).Where("user_id = ? and corporation_name = ?", credentials.UserID, credentials.CorporationName).Delete(); err != nil {
		return fmt.Errorf("error when getting corporation credentials for userID %d: %w", userID, err)
	}

	return nil
}

func (s *service) DecryptCredentials(credentials *customer.CorporationCredentials) (*customer.CorporationCredentials, error) {
	// decrypt credentials
	var err error
	secretKey := auth.BuildAESKey(credentials.UserID, credentials.CorporationName, s.aesSecret)
	if credentials.Login, err = auth.DecryptString(credentials.Login, secretKey); err != nil {
		return nil, fmt.Errorf("error when decrypting corporation credentials: %w", err)
	}

	if credentials.Password, err = auth.DecryptString(credentials.Password, secretKey); err != nil {
		return nil, fmt.Errorf("error when decrypting corporation credentials: %w", err)
	}

	return credentials, nil
}

func (s *service) validateCredentials(credentials customer.CorporationCredentials) error {
	client, err := s.clientProvider.Get(credentials.CorporationName)
	if err != nil {
		return err
	}
	if err := client().Login(credentials.Login, credentials.Password); err != nil {
		return fmt.Errorf("error when authenticating to %s with given credentials: %w", credentials.CorporationName, err)
	}

	return nil
}
