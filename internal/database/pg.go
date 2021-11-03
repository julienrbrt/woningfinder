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
func NewDBClient(logger *logging.Logger, debug bool, host, port, name, user, password string) (DBClient, error) {
	dbLogger := &dbLogger{logger}

	// set logger
	pg.SetLogger(dbLogger)

	// in debug mode ssl is not required (debug mode should only be ran locally)
	sslmode := "require"
	if debug {
		sslmode = "disable"
	}

	// connect to database
	opt, err := pg.ParseURL(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&", user, password, host, port, name, sslmode))
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
