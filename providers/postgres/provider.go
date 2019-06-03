package postgres

import (
	"time"

	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

type (
	Provider struct {
		db     *pg.DB
		logger *zap.Logger
	}
)

func NewProvider(postgresURL string, logger *zap.Logger) (*Provider, error) {
	opt, err := pg.ParseURL(postgresURL)
	if err != nil {
		return nil, err
	}

	opt.PoolSize = 10
	db := pg.Connect(opt)

	return &Provider{db, logger}, nil
}

func (provider *Provider) GetDatabase() Database {
	return &postgresDB{provider.db}
}

func (provider *Provider) NotifyFail(ce chan error) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			<-ticker.C

			_, err := provider.db.Exec("SELECT 1")
			if err != nil {
				provider.logger.Error("health check to postgres failed", zap.Error(err))
				ce <- err
				return
			}
		}
	}()
}

func (provider *Provider) Close() {
	if err := provider.db.Close(); err != nil {
		provider.logger.Error("unable to close connection properly", zap.Error(err))
	}
}

func (provider *Provider) GetClient() *pg.DB {
	return provider.db
}
