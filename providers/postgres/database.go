package postgres

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type (
	Database interface {
		orm.DB

		Begin() (Database, error)
		Rollback() error
		Commit() error
	}

	postgresDB struct {
		*pg.DB
	}

	postgresTX struct {
		*pg.Tx
	}
)

func (postgres *postgresDB) Begin() (Database, error) {
	tx, err := postgres.DB.Begin()
	if err != nil {
		return nil, err
	}

	return &postgresTX{tx}, nil
}

func (postgres *postgresDB) Rollback() error {
	return nil
}

func (postgres *postgresDB) Commit() error {
	return nil
}

func (postgres *postgresTX) Begin() (Database, error) {
	return postgres, nil
}
