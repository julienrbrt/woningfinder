package database

import (
	"context"
	"fmt"

	"github.com/woningfinder/woningfinder/pkg/logging"

	"github.com/go-pg/pg/v10"
)

type DBClient interface {
	Conn() *pg.DB
}

type dbClient struct {
	logger *logging.Logger
	db     *pg.DB
}

// NewDBClient create a connection to postgresql database
func NewDBClient(logger *logging.Logger, host, port, name, user, password string) (DBClient, error) {
	// set logger
	pg.SetLogger(logger)

	// connect to database
	opt, err := pg.ParseURL(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require&", user, password, host, port, name))
	if err != nil {
		return nil, err
	}
	db := pg.Connect(opt)
	// log each query
	db.AddQueryHook(dbLogger{logger: logger})

	// check connection
	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		panic(err)
	} else {
		logger.Info("successfully connected to postgresql 🎉")
	}

	return &dbClient{
		logger: logger,
		db:     db,
	}, nil
}

func (db *dbClient) Conn() *pg.DB {
	return db.db
}
