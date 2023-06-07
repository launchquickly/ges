// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package sqlc

import (
	"encoding/json"

	"github.com/google/uuid"
)

type EventStoreAggregate struct {
	ID            uuid.UUID
	Version       int32
	AggregateType string
}

type EventStoreAggregateSnapshot struct {
	AggregateID uuid.UUID
	Version     int32
	Data        json.RawMessage
}

type EventStoreEvent struct {
	ID            int64
	TransactionID interface{}
	AggregateID   uuid.UUID
	Version       int32
	Data          json.RawMessage
}
