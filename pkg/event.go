package es

import (
	"reflect"
)

// Event is a domain event marker
type Event interface {
	IsEvent()
}

func EventName(event Event) string {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}
