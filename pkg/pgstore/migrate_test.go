package pgstore_test

import (
	"github.com/launchquickly/ges/pkg/pgstore"
	"github.com/launchquickly/ges/pkg/pgstore/test"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	expectedUpVersion uint = 2
	sourceURL              = "file://./migration"
)

func TestMigrateDown(t *testing.T) {
	database := test.NewDatabase(t)
	c := database.PasswordConfig(t)

	version, err := pgstore.MigrateUp(sourceURL, c, true)
	require.NoError(t, err)
	require.Equal(t, expectedUpVersion, version)

	version, err = pgstore.MigrateDown(sourceURL, c, true)
	require.NoError(t, err)
	require.Equal(t, uint(0), version)
}

func TestMigrateUp(t *testing.T) {
	database := test.NewDatabase(t)
	c := database.PasswordConfig(t)

	version, err := pgstore.MigrateUp(sourceURL, c, true)
	require.NoError(t, err)
	require.Equal(t, expectedUpVersion, version)
}
