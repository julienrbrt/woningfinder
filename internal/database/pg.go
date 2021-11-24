package database

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/julienrbrt/woningfinder/pkg/logging"
)

type DBClient interface {
	Conn() *pg.DB
}

type dbClient struct {
	logger *logging.Logger
	db     *pg.DB
}

// NewDBClient create a connection to postgresql database
func NewDBClient(logger *logging.Logger, debug bool, databaseURL string) (DBClient, error) {
	dbLogger := &dbLogger{logger}

	// set logger
	pg.SetLogger(dbLogger)

	// connect to database
	opt, err := pg.ParseURL(databaseURL)
	if err != nil {
		return nil, err
	}
	db := pg.Connect(opt)

	// log each query on debug mode
	if debug {
		db.AddQueryHook(dbLogger)
	}

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
