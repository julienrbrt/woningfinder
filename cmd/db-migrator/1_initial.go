package main

import (
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10/orm"
	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
	"github.com/julienrbrt/woningfinder/internal/customer"
	"github.com/julienrbrt/woningfinder/internal/services/matcher"
	"github.com/julienrbrt/woningfinder/pkg/mapbox"
)

func init() {
	// models
	models := []interface{}{
		(*city.City)(nil),
		(*mapbox.AddressCityDistrict)(nil),
		(*matcher.MatcherCounter)(nil),
		(*customer.User)(nil),
		(*corporation.Corporation)(nil),
		(*corporation.CorporationCity)(nil),
		(*customer.HousingPreferences)(nil),
		(*customer.HousingPreferencesHousingType)(nil),
		(*customer.HousingPreferencesCity)(nil),
		(*customer.HousingPreferencesCityDistrict)(nil),
		(*customer.HousingPreferencesMatch)(nil),
		(*customer.CorporationCredentials)(nil),
		(*customer.ReminderCounter)(nil),
		(*customer.WaitingList)(nil),
	}

	migrations.MustRegisterTx(func(db migrations.DB) error {
		for _, model := range models {
			err := db.Model(model).CreateTable(&orm.CreateTableOptions{
				IfNotExists: true,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
}
