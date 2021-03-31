package main

import (
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10/orm"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

func init() {
	// register m2m models with go-pg
	bootstrap.RegisterModel()

	// models
	models := []interface{}{
		(*entity.CorporationCity)(nil),
		(*entity.HousingPreferencesHousingType)(nil),
		(*entity.HousingPreferencesCity)(nil),

		(*entity.Corporation)(nil),
		(*entity.HousingType)(nil),
		(*entity.City)(nil),
		(*entity.User)(nil),
		(*entity.UserPlan)(nil),
		(*entity.HousingPreferences)(nil),
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
