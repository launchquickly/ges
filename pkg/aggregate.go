package es

import (
	"reflect"
)

// Aggregate defines the base aggregate operations.
type Aggregate interface {
	ID() ID
}

func AggregateType(aggregate Aggregate) string {
	t := reflect.TypeOf(aggregate)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}
