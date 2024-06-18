package repository

import (
	"context"
	"github.com/go-pg/pg/v10"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Config struct {
	URL       string `env:"COCKROACH_CLIENT_DATABASE_URL_GOLANG"`
	DebugMode bool   `env:"COCKROACH_CLIENT_DB_DEBUG" envDefault:"false"`
}

func NewDatabase(cfg *Config) (*pg.DB, error) {
	opt, err := pg.ParseURL(cfg.URL)
	if err != nil {
		return nil, err
	}
	db := pg.Connect(opt)
	if cfg.DebugMode {
		db.AddQueryHook(dbLogger{})
	}
	if err = db.Ping(context.Background()); err != nil {
		return nil, err
	}
	return db, nil
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	if query, err := q.FormattedQuery(); err != nil {
		zap.S().Infof("[SQL]: %s", err.Error())
	} else {
		zap.S().Infof("[SQL]: %s", string(query))
	}

	return nil
}
