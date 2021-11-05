package main

import (
	"github.com/go-pg/migrations/v8"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/woningnet"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		for _, corp := range []corporation.Corporation{
			woningnet.AlmereInfo,
			woningnet.WoonkeusInfo,
			woningnet.EemvalleiInfo,
			woningnet.WoonserviceInfo,
			woningnet.MercatusInfo,
			woningnet.MiddenHollandInfo,
			woningnet.BovenGroningenInfo,
			woningnet.GooiVechtstreekInfo,
			woningnet.GroningenInfo,
			woningnet.HuiswaartsInfo,
			woningnet.WoongaardInfo,
		} {
			if err := corporationService.CreateOrUpdateCorporation(corp); err != nil {
				return err
			}
		}

		return nil
	})
}
