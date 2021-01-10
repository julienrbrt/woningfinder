package entity

import (
	"gorm.io/gorm"
)

// CorporationCredentials holds the user credentials to login to an housing corporation
type CorporationCredentials struct {
	gorm.Model
	UserID          int    `gorm:"primaryKey"`
	CorporationName string `gorm:"primaryKey"`
	CorporationURL  string `gorm:"primaryKey"`
	Corporation     Corporation
	Login           string
	Password        string
}
