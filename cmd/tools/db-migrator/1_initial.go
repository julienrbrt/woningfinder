package main

import (
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10/orm"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

func init() {
	// models
	models := []interface{}{
		(*entity.City)(nil),
		(*entity.User)(nil),
		(*entity.UserPlan)(nil),
		(*entity.Corporation)(nil),
		(*entity.CorporationCity)(nil),
		(*entity.HousingPreferences)(nil),
		(*entity.HousingPreferencesHousingType)(nil),
		(*entity.HousingPreferencesCity)(nil),
		(*entity.HousingPreferencesCityDistrict)(nil),
		(*entity.HousingPreferencesMatch)(nil),
		(*entity.CorporationCredentials)(nil),
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
