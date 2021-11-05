package main

import (
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10/orm"
	"github.com/woningfinder/woningfinder/internal/customer"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		return db.Model((*customer.WaitingList)(nil)).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
	})
}
