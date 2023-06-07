package test

import (
	"context"
	"database/sql"
	store "github.com/launchquickly/ges/pkg"
	"github.com/launchquickly/ges/pkg/pgstore/sqlc"
	"github.com/stretchr/testify/require"
	"testing"
)

func InsertRecord(t *testing.T, dataSource string, aggregate store.Aggregate, params sqlc.AppendEventParams) {
	t.Helper()

	db, err := sql.Open("postgres", dataSource)
	require.NoError(t, err)

	tx, err := db.Begin()
	require.NoError(t, err)
	defer func() {
		err = tx.Rollback()
	}()

	queries := sqlc.New(db)
	qtx := queries.WithTx(tx)

	ctxt := context.Background()

	if params.Version == 1 {
		qtx.CreateAggregate(ctxt, sqlc.CreateAggregateParams{
			AggregateID:   aggregate.ID(),
			AggregateType: store.AggregateType(aggregate),
		})
	}

	if _, err := qtx.AppendEvent(ctxt, params); err != nil {
		require.NoError(t, err)
	}
	tx.Commit()
}

func InsertRecords(t *testing.T, dataSource string, aggregate store.Aggregate, records ...store.Record) {
	t.Helper()

	for _, r := range records {
		InsertRecord(t, dataSource, aggregate, sqlc.AppendEventParams{
			AggregateID: aggregate.ID(),
			Data:        r.Data,
			Version:     r.Version,
		})
	}
}

func LoadStream(t *testing.T, s store.Store, aggregateID store.ID, from int32, to int32) store.Stream {
	t.Helper()

	stream, err := s.Load(context.Background(), aggregateID, from, to)
	require.NoError(t, err)
	return stream
}
