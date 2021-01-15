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
		(*entity.CorporationSelectionMethod)(nil),
		(*entity.HousingPreferencesHousingType)(nil),
		(*entity.HousingPreferencesCity)(nil),
		(*entity.HousingPreferencesCityDistrict)(nil),

		(*entity.Corporation)(nil),
		(*entity.SelectionMethod)(nil),
		(*entity.HousingType)(nil),
		(*entity.City)(nil),
		(*entity.CityDistrict)(nil),
		(*entity.User)(nil),
		(*entity.Tier)(nil),
		(*entity.HousingPreferences)(nil),
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
