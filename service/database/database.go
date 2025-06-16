package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Database interface {
	Ping() error
	Close() error
}

type databaseImpl struct {
	connection *sql.DB
}

func New(db *sql.DB, migrationsPath string) (Database, error) {
	if db == nil {
		return nil, errors.New("database connection is required")
	}

	if migrationsPath == "" {
		return nil, errors.New("migrations path is required")
	}

	migrations, err := os.ReadDir(migrationsPath)

	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, migration := range migrations {
		if strings.HasSuffix(migration.Name(), ".sql") {
			migrationPath := filepath.Join(migrationsPath, migration.Name())

			migrationFile, err := os.Open(migrationPath)

			if err != nil {
				return nil, fmt.Errorf("failed to open migration file %s: %w", migration.Name(), err)
			}

			defer migrationFile.Close()

			migrationSQL, err := io.ReadAll(migrationFile)

			if err != nil {
				return nil, fmt.Errorf("failed to read migration file %s: %w", migration.Name(), err)
			}

			_, err = db.Exec(string(migrationSQL))

			if err != nil {
				return nil, fmt.Errorf("failed to execute migration %s: %w", migration.Name(), err)
			}
		}
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
