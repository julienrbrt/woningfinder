package user

import (
	"errors"
	"fmt"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/util"
	"gorm.io/gorm"
)

// CorporationCredentials holds the user credentials to login to an housing corporation
type CorporationCredentials struct {
	gorm.Model
	UserID          int    `gorm:"primaryKey"`
	CorporationName string `gorm:"primaryKey"`
	CorporationURL  string `gorm:"primaryKey"`
	Corporation     corporation.Corporation
	Login           string
	Password        string
}

func (s *service) CreateCorporationCredentials(u *User, credentials CorporationCredentials) error {
	if credentials.Corporation.Name == "" || credentials.Login == "" || credentials.Password == "" {
		return fmt.Errorf("error login or password cannot be empty when adding credentials")
	}

	// check credentials validity
	client, err := s.clientProvider.Get(credentials.Corporation)
	if err != nil {
		return err
	}
	if err := client.Login(credentials.Login, credentials.Password); err != nil {
		return fmt.Errorf("error when authenticating to %s with given credentials: %w", credentials.Corporation.Name, err)
	}

	// encrypt credentials
	credentials.Login, err = util.AESEncrypt(credentials.Login, s.aesSecret)
	if err != nil {
		return fmt.Errorf("error when encrypting corporation credentials: %w", err)
	}

	credentials.Password, err = util.AESEncrypt(credentials.Password, s.aesSecret)
	if err != nil {
		return fmt.Errorf("error when encrypting corporation credentials: %w", err)
	}

	// check if already existing
	fetchCredentials, err := s.GetCorporationCredentials(u, credentials.Corporation)
	if err != nil && !errors.Is(err, errCorporationCredentialsNotFound) {
		return fmt.Errorf("error when checking if credentials already exists: %w", err)
	}

	// store credentials
	if err != nil { // store unexisting credentials
		if err := s.dbClient.Conn().Model(u).Association("CorporationCredentials").Append(&credentials); err != nil {
			return fmt.Errorf("error when creating corporation credentials: %w", err)
		}
	} else { // update existing credentials
		if err := s.dbClient.Conn().Model(&fetchCredentials).Updates(&credentials).Error; err != nil {
			return fmt.Errorf("error when updating corporation credentials: %w", err)
		}
	}

	return nil
}

func (s *service) GetCorporationCredentials(u *User, corporation corporation.Corporation) (*CorporationCredentials, error) {
	query := CorporationCredentials{
		UserID:          int(u.ID),
		CorporationName: corporation.Name,
		CorporationURL:  corporation.URL,
	}

	// get corporation credentials
	var credentials CorporationCredentials
	if err := s.dbClient.Conn().Where(query).Find(&credentials).Error; err != nil {
		return nil, fmt.Errorf("error when getting corporation credentials for user %s: %w", u.Email, err)
	}

	if credentials.Login == "" || credentials.Password == "" {
		return nil, errCorporationCredentialsNotFound
	}

	// decrypt credentials
	var err error
	credentials.Login, err = util.AESDecrypt(credentials.Login, s.aesSecret)
	if err != nil {
		return nil, fmt.Errorf("error when decrypting corporation credentials: %w", err)
	}

	credentials.Password, err = util.AESDecrypt(credentials.Password, s.aesSecret)
	if err != nil {
		return nil, fmt.Errorf("error when decrypting corporation credentials: %w", err)
	}

	return &credentials, nil
}

func (s *service) DeleteCorporationCredentials(u *User, corporation corporation.Corporation) error {
	credentials, err := s.GetCorporationCredentials(u, corporation)
	if err != nil {
		return fmt.Errorf("error when deleting corporation credentials: %w", err)
	}

	// delete permanently
	credentials.Login = ""
	credentials.Password = ""
	if err = s.dbClient.Conn().Unscoped().Delete(credentials).Error; err != nil {
		return fmt.Errorf("error when deleting corporation credentials: %w", err)
	}

	return nil
}

func (s *service) decryptCredentials(credentials CorporationCredentials) (CorporationCredentials, error) {
	// decrypt credentials
	var err error
	credentials.Login, err = util.AESDecrypt(credentials.Login, s.aesSecret)
	if err != nil {
		return CorporationCredentials{}, fmt.Errorf("error when decrypting corporation credentials: %w", err)
	}

	credentials.Password, err = util.AESDecrypt(credentials.Password, s.aesSecret)
	if err != nil {
		return CorporationCredentials{}, fmt.Errorf("error when decrypting corporation credentials: %w", err)
	}

	return credentials, nil
}
