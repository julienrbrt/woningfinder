package database

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

type dbLogger struct {
	logger *logging.Logger
}

func (db dbLogger) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (db dbLogger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	query, err := q.FormattedQuery()
	if err != nil {
		return err
	}

	db.logger.Sugar().Debugf("go-pg query log: %s", string(query))
	return nil
}
