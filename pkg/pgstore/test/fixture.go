package test

import (
	"encoding/json"
	"github.com/google/uuid"
	store "github.com/launchquickly/ges/pkg"
)

var (
	SingleRecord = store.Record{
		Version: 1,
		Data:    json.RawMessage(`{ "type": "single", "field": "value" }`),
	}

	OneRecord = store.Record{
		Version: 1,
		Data:    json.RawMessage(`{ "type": "one", "field": "value1" }`),
	}
	TwoRecord = store.Record{
		Version: 2,
		Data:    json.RawMessage(`{ "type": "two", "field": "value2" }`),
	}
	ThreeRecord = store.Record{
		Version: 3,
		Data:    json.RawMessage(`{ "type": "three", "field": "value3" }`),
	}
	FourRecord = store.Record{
		Version: 4,
		Data:    json.RawMessage(`{ "type": "four", "field": "value4" }`),
	}
	FiveRecord = store.Record{
		Version: 5,
		Data:    json.RawMessage(`{ "type": "five", "field": "value5" }`),
	}
	SixRecord = store.Record{
		Version: 6,
		Data:    json.RawMessage(`{ "type": "six", "field": "value6" }`),
	}
)

type Fixture struct {
	Aggregate      store.Aggregate
	Records        []store.Record
	FromVersion    int32
	ToVersion      int32
	ExpectedLength int
	buildCalled    bool
	setID          bool
}

func NewFixture() *Fixture {
	return &Fixture{}
}

func (f *Fixture) Build() *Fixture {
	if !f.setID {
		f.Aggregate = AStub{
			id: uuid.New(),
		}
		f.setID = true
	}
	f.buildCalled = true
	return f
}

func (f *Fixture) From(from int32) *Fixture {
	f.FromVersion = from
	return f
}

func (f *Fixture) To(to int32) *Fixture {
	f.ToVersion = to
	return f
}

func (f *Fixture) For(aggregateID store.ID) *Fixture {
	f.Aggregate = AStub{
		id: aggregateID,
	}
	f.setID = true
	return f
}

func (f *Fixture) With(records ...store.Record) *Fixture {
	f.Records = records
	return f
}

func (f *Fixture) OfLength(expectedLength int) *Fixture {
	f.ExpectedLength = expectedLength
	return f
}

type AStub struct {
	id store.ID
}

func (a AStub) ID() store.ID {
	return a.id
}
