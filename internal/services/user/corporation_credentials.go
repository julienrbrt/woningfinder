package user

import (
	"errors"
	"fmt"

	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/entity"
	"github.com/woningfinder/woningfinder/internal/services"
)

func (s *service) CreateCorporationCredentials(userID uint, credentials entity.CorporationCredentials) error {
	if credentials.CorporationName == "" {
		return errors.New("error when creating corporation credentials: corporation invalid")
	}

	if credentials.Login == "" || credentials.Password == "" {
		return errors.New("error when creating corporation credentials: login or password missing")
	}

	// check credentials validity
	if err := s.ValidateCredentials(credentials); err != nil {
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

func (s *service) GetCorporationCredentials(userID uint, corporation entity.Corporation) (*entity.CorporationCredentials, error) {
	credentials := entity.CorporationCredentials{
		UserID:          userID,
		CorporationName: corporation.Name,
	}

	// get corporation credentials
	if err := s.dbClient.Conn().Model(&credentials).Where("user_id = ? and corporation_name ILIKE ?", credentials.UserID, credentials.CorporationName).Select(); err != nil {
		return nil, fmt.Errorf("error when getting corporation credentials for userID %d: %w", userID, err)
	}

	return &credentials, nil
}

func (s *service) GetAllCorporationCredentials(corporation entity.Corporation) ([]entity.CorporationCredentials, error) {
	credentials := []entity.CorporationCredentials{}
	if err := s.dbClient.Conn().
		Model(&credentials).
		Where("corporation_name ILIKE ?", corporation.Name).
		Order("created_at ASC"). // people having registered their credentials for the longer get reaction priority (see documentation)
		Select(); err != nil {
		return nil, fmt.Errorf("error getting user credentials: %w", err)
	}

	// no users found
	if len(credentials) == 0 {
		return nil, services.ErrNoMatchFound
	}

	return credentials, nil
}

func (s *service) DeleteCorporationCredentials(userID uint, corporation entity.Corporation) error {
	credentials, err := s.GetCorporationCredentials(userID, corporation)
	if err != nil {
		return fmt.Errorf("error when getting corporation credentials: %w", err)
	}

	// delete permanently
	credentials.Login = ""
	credentials.Password = ""

	// TODO to implement
	// delete corporations credentials
	panic("not implemented")
}

func (s *service) ValidateCredentials(credentials entity.CorporationCredentials) error {
	client, err := s.clientProvider.GetByName(entity.Corporation{Name: credentials.CorporationName})
	if err != nil {
		return err
	}
	if err := client.Login(credentials.Login, credentials.Password); err != nil {
		return fmt.Errorf("error when authenticating to %s with given credentials: %w", credentials.CorporationName, err)
	}

	return nil
}

func (s *service) DecryptCredentials(credentials *entity.CorporationCredentials) (*entity.CorporationCredentials, error) {
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
