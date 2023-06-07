package pgstore

import (
	"context"
	"database/sql"
	"fmt"
	store "github.com/launchquickly/ges/pkg"
	"github.com/launchquickly/ges/pkg/pgstore/sqlc"
	_ "github.com/lib/pq"
	"math"
)

func New(ds string) *PostgresStore {
	return &PostgresStore{
		dataSource: ds,
	}
}

type PostgresStore struct {
	dataSource string
}

// Load the stream of events up to the version specified.
// When toVersion is 0, all events will be loaded.
// To start at the beginning, fromVersion should be set to 0
func (s *PostgresStore) Load(ctx context.Context, aggregateID store.ID, fromVersion, toVersion int32) (store.Stream, error) {

	if toVersion == 0 {
		toVersion = math.MaxInt32
	}

	db, err := sql.Open(driverName, s.dataSource)
	if err != nil {
		return nil, err
	}

	queries := sqlc.New(db)

	// load stream
	params := sqlc.LoadStreamParams{
		AggregateID: aggregateID,
		FromVersion: fromVersion,
		ToVersion:   toVersion,
	}
	rows, err := queries.LoadStream(ctx, params)
	if err != nil {
		return nil, err
	}

	stream := store.Stream{}
	for _, row := range rows {
		stream = append(stream, store.Record{
			Version: row.Version,
			Data:    row.Data,
		})
	}

	return stream, nil
}

// Save the provided serialized records to the store
func (s *PostgresStore) Save(ctx context.Context, aggregate store.Aggregate, records ...store.Record) error {

	if len(records) == 0 {
		return nil
	}

	err := checkRecordSequence(records)
	if err != nil {
		return err
	}

	db, err := sql.Open(driverName, s.dataSource)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		err = tx.Rollback()
	}()

	queries := sqlc.New(db)
	qtx := queries.WithTx(tx)

	if records[0].Version == 1 {
		err = qtx.CreateAggregate(ctx, sqlc.CreateAggregateParams{
			AggregateID:   aggregate.ID(),
			AggregateType: store.AggregateType(aggregate),
		})
		if err != nil {
			return err
		}
	}

	ev := records[0].Version - 1
	ru, err := qtx.UpdateAggregate(ctx, sqlc.UpdateAggregateParams{
		AggregateID:     aggregate.ID(),
		ExpectedVersion: ev,
		NewVersion:      records[len(records)-1].Version,
	})
	if err != nil {
		return err
	}
	if ru != 1 {
		return fmt.Errorf("optimistic concurrency version check detected conflict. actual version does not match expected version: %d", ev)
	}

	for _, record := range records {
		params := sqlc.AppendEventParams{
			AggregateID: aggregate.ID(),
			Data:        record.Data,
			Version:     record.Version,
		}
		if _, err := qtx.AppendEvent(ctx, params); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func checkRecordSequence(records []store.Record) error {
	if len(records) < 2 {
		return nil
	}
	pv := records[0].Version
	for i := 1; i < len(records); i++ {
		cv := records[i].Version
		if cv != pv+1 {
			return fmt.Errorf("sequence check detected issue. version: %d is out of sequence", cv)
		}
		pv = cv
	}
	return nil
}
