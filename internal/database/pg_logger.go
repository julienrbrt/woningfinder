package database

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"go.uber.org/zap"
)

type dbLogger struct {
	*logging.Logger
}

// Printf prints a log as info
func (l *dbLogger) Printf(_ context.Context, template string, args ...interface{}) {
	l.Sugar().Infof(template, args)
}

func (*dbLogger) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (l *dbLogger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	query, err := q.FormattedQuery()
	if err != nil {
		return err
	}

	l.Debug("go-pg query log", zap.String("query", string(query)))
	return nil
}
