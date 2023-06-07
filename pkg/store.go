package es

import (
	"context"
	"encoding/json"
)

// Record provides the serialized representation of the event
type Record struct {
	// Version contains the version associated with the serialized event
	Version int32

	// Data contains the event in serialized form
	Data json.RawMessage
}

// Stream represents
type Stream []Record

// Store provides an abstraction for the Repository to save data
type Store interface {
	// Save the provided serialized records to the store
	Save(ctx context.Context, aggregate Aggregate, records ...Record) error

	// Load the stream of events up to the version specified.
	// When toVersion is 0, all events will be loaded.
	// To start at the beginning, fromVersion should be set to 0
	// fromVersion and toVersion values are inclusive
	Load(ctx context.Context, aggregateID ID, fromVersion, toVersion int32) (Stream, error)
}
