package pgstore

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // blank import for test - this comment is to appease revive
)

// MigrateDown carries out database migration applying down migrations.
// Migrations are loaded from sourceURL directory on the host.
//
// MigrateDown should be expected to result in a version of 0 as it reverts all previous Ups.
//
// * sourceURL - location of the migrations
// * c         - configuration for the database to carry out migration on
// * l         - whether executing locally or on AWS
//
// returns currently active migration version.
func MigrateDown(sourceURL string, c PasswordConfig, l bool) (uint, error) {
	m, err := NewMigrate(sourceURL, c, l)
	if err != nil {
		return 0, err
	}

	if err = m.Down(); err != nil && err != migrate.ErrNoChange {
		return 0, err
	}

	return currentActiveVersion(m)
}

// MigrateUp carries out database migration applying up migrations.
// Migrations are loaded from sourceURL directory on the host.
//
// * sourceURL - location of the migrations
// * c         - configuration for the database to carry out migration on
// * l         - whether executing locally or on AWS
//
// returns currently active migration version.
func MigrateUp(sourceURL string, c PasswordConfig, l bool) (uint, error) {
	m, err := NewMigrate(sourceURL, c, l)
	if err != nil {
		return 0, err
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return 0, err
	}

	return currentActiveVersion(m)
}

func NewMigrate(sourceURL string, c PasswordConfig, l bool) (*migrate.Migrate, error) {
	passwordDB, err := NewPasswordDB(c, l)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(passwordDB, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(sourceURL, c.Name, driver)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func currentActiveVersion(m *migrate.Migrate) (uint, error) {
	version, _, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return 0, err
	}
	return version, nil
}
