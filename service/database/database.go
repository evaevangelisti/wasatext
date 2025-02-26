package database

import (
	"database/sql"
	"errors"
)

type Database interface {
	Ping() error
	Close() error
}

type databaseImpl struct {
	connection *sql.DB
}

func New(db *sql.DB) (Database, error) {
	if db == nil {
		return nil, errors.New("database connection is required")
	}

	return &databaseImpl{
		connection: db,
	}, nil
}

func (db *databaseImpl) Ping() error {
	return db.connection.Ping()
}

func (db *databaseImpl) Close() error {
	return db.connection.Close()
}
