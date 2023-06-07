package test

import (
	"context"
	"fmt"
	"github.com/launchquickly/ges/pkg/pgstore"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

const (
	dbName     = "go_event_store"
	dbPassword = "secret"
	dbUser     = "postgres"
)

type Database struct {
	instance testcontainers.Container
}

func NewDatabase(t *testing.T) *Database {
	t.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15.3",
		ExposedPorts: []string{"5432/tcp"},
		HostConfigModifier: func(config *container.HostConfig) {
			config.AutoRemove = true
		},
		Env: map[string]string{
			"POSTGRES_USER":             dbUser,
			"POSTGRES_PASSWORD":         dbPassword,
			"POSTGRES_DB":               dbName,
			"POSTGRES_HOST_AUTH_METHOD": "trust",
			"PGTZ":                      "0",
		},
		WaitingFor: NewPostgresWaitStrategy(),
	}
	genericContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	require.NoError(t, err)

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		require.NoError(t, genericContainer.Terminate(ctx))
	})

	return &Database{
		instance: genericContainer,
	}
}

func (tdb *Database) PasswordConfig(t *testing.T) pgstore.PasswordConfig {
	t.Helper()
	return pgstore.PasswordConfig{
		Host:     "localhost",
		Name:     dbName,
		Port:     tdb.Port(t),
		Username: dbUser,
		Password: dbPassword,
	}
}

func (tdb *Database) Port(t *testing.T) int {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	p, err := tdb.instance.MappedPort(ctx, "5432")
	require.NoError(t, err)
	return p.Int()
}

func (tdb *Database) ConnectionString(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("postgres://%s:%s@127.0.0.1:%d/%s?sslmode=disable", dbUser, dbPassword, tdb.Port(t), dbName)
}
