package pgstore_test

import (
	"context"
	store "github.com/launchquickly/ges/pkg"
	"github.com/launchquickly/ges/pkg/pgstore"
	"github.com/launchquickly/ges/pkg/pgstore/test"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	database *test.Database
	ps       *pgstore.PostgresStore
)

func TestPostgresStore_Load(t *testing.T) {

	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name    string
		fixture *test.Fixture
	}{
		{
			name:    "no records returns empty stream",
			fixture: test.NewFixture().OfLength(0).Build(),
		},
		{
			name:    "single record",
			fixture: test.NewFixture().With(test.SingleRecord).OfLength(1).Build(),
		},
		{
			name: "multiple records",
			fixture: test.NewFixture().With(test.OneRecord, test.TwoRecord, test.ThreeRecord, test.FourRecord).
				OfLength(4).Build(),
		},
		{
			name: "from first record",
			fixture: test.NewFixture().With(test.OneRecord, test.TwoRecord, test.ThreeRecord).From(1).
				OfLength(3).Build(),
		},
		{
			name: "from third record",
			fixture: test.NewFixture().
				With(test.OneRecord, test.TwoRecord, test.ThreeRecord, test.FourRecord, test.FiveRecord).From(3).
				OfLength(3).Build(),
		},
		{
			name: "from last record",
			fixture: test.NewFixture().With(test.OneRecord, test.TwoRecord).From(2).OfLength(1).
				Build(),
		},
		{
			name:    "to first record",
			fixture: test.NewFixture().With(test.OneRecord, test.TwoRecord).To(1).OfLength(1).Build(),
		},
		{
			name: "to fourth record",
			fixture: test.NewFixture().
				With(test.OneRecord, test.TwoRecord, test.ThreeRecord, test.FourRecord, test.FiveRecord).To(4).
				OfLength(4).Build(),
		},
		{
			name: "to last record",
			fixture: test.NewFixture().With(test.OneRecord, test.TwoRecord, test.ThreeRecord).To(3).
				OfLength(3).Build(),
		},
		{
			name: "from second to fourth record",
			fixture: test.NewFixture().
				With(test.OneRecord, test.TwoRecord, test.ThreeRecord, test.FourRecord, test.FiveRecord).From(2).To(4).
				OfLength(3).Build(),
		},
		{
			name: "from second to second record",
			fixture: test.NewFixture().
				With(test.OneRecord, test.TwoRecord, test.ThreeRecord).From(2).To(2).OfLength(1).Build(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			teardownTest := setupTest(t)
			defer teardownTest(t)

			// given
			fixture := tt.fixture
			aggregate := fixture.Aggregate
			records := fixture.Records
			from := fixture.FromVersion
			to := fixture.ToVersion
			expectedLength := fixture.ExpectedLength

			// save records to later load
			test.InsertRecords(t, database.ConnectionString(t), aggregate, records...)

			// when
			stream, err := ps.Load(context.Background(), aggregate.ID(), from, to)
			require.NoError(t, err)

			// then
			test.AssertStreamLength(t, stream, expectedLength)
			test.AssertStreamContents(t, fixture, stream)
		})
	}
}

func TestPostgresStore_Save(t *testing.T) {

	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name    string
		records []store.Record
		fixture *test.Fixture
	}{
		{
			name:    "no records does not fail",
			fixture: test.NewFixture().OfLength(0).Build(),
		},
		{
			name:    "single record",
			fixture: test.NewFixture().With(test.SingleRecord).OfLength(1).Build(),
		},
		{
			name: "multiple records",
			fixture: test.NewFixture().With(test.SingleRecord, test.TwoRecord, test.ThreeRecord).OfLength(3).
				Build(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			teardownTest := setupTest(t)
			defer teardownTest(t)

			// given
			fixture := tt.fixture
			aggregate := tt.fixture.Aggregate
			records := fixture.Records
			from := fixture.FromVersion
			to := fixture.ToVersion
			expectedLength := fixture.ExpectedLength

			// when
			err := ps.Save(context.Background(), aggregate, records...)
			require.NoError(t, err)

			// then
			stream := test.LoadStream(t, ps, aggregate.ID(), from, to)

			test.AssertStreamLength(t, stream, expectedLength)
			test.AssertStreamContents(t, fixture, stream)
		})
	}
}

func TestPostgresStore_Save_Append(t *testing.T) {

	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name string
		// position we want to append from
		appendOffset int
		fixture      *test.Fixture
	}{
		{
			name:         "no records does not fail",
			appendOffset: 3,
			fixture: test.NewFixture().With(test.OneRecord, test.TwoRecord, test.ThreeRecord).OfLength(3).
				Build(),
		},
		{
			name:         "single record",
			appendOffset: 2,
			fixture: test.NewFixture().With(test.OneRecord, test.TwoRecord, test.ThreeRecord).OfLength(3).
				Build(),
		},
		{
			name:         "multiple records",
			appendOffset: 1,
			fixture: test.NewFixture().With(test.OneRecord, test.TwoRecord, test.ThreeRecord, test.FourRecord, test.FiveRecord,
				test.SixRecord).OfLength(6).Build(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			teardownTest := setupTest(t)
			defer teardownTest(t)

			// given
			fixture := tt.fixture
			aggregate := fixture.Aggregate
			records := fixture.Records
			from := fixture.FromVersion
			to := fixture.ToVersion
			appendOffset := tt.appendOffset
			expectedLength := fixture.ExpectedLength

			// save some records so we can append to them
			err := ps.Save(context.Background(), aggregate, records[:appendOffset]...)
			require.NoError(t, err)

			// when
			err = ps.Save(context.Background(), aggregate, records[appendOffset:]...)
			require.NoError(t, err)

			// then
			stream, err := ps.Load(context.Background(), aggregate.ID(), from, to)
			require.NoError(t, err)

			test.AssertStreamLength(t, stream, expectedLength)
			test.AssertStreamContents(t, fixture, stream)
		})
	}
}

func TestPostgresStore_Save_Errors(t *testing.T) {

	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name    string
		records []store.Record
		fixture *test.Fixture
	}{
		{
			name:    "first record saved must be version 1",
			fixture: test.NewFixture().With(test.TwoRecord).OfLength(1).Build(),
		},
		{
			name:    "version numbers must be sequential",
			fixture: test.NewFixture().With(test.SingleRecord, test.ThreeRecord).OfLength(3).Build(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			teardownTest := setupTest(t)
			defer teardownTest(t)

			// given
			fixture := tt.fixture
			aggregate := fixture.Aggregate
			records := fixture.Records

			// when
			err := ps.Save(context.Background(), aggregate, records...)
			require.Error(t, err)

			// check details of err
			// TODO
		})
	}
}

func setupSuite(t *testing.T) func(t *testing.T) {
	t.Helper()
	database = test.NewDatabase(t)
	c := database.PasswordConfig(t)

	_, err := pgstore.MigrateUp(sourceURL, c, true)
	require.NoError(t, err)

	// Return a function to teardown the test
	return func(t *testing.T) {
		t.Helper()
	}
}

func setupTest(t *testing.T) func(t *testing.T) {
	t.Helper()
	ps = pgstore.New(database.ConnectionString(t))

	// Return a function to teardown the test
	return func(t *testing.T) {
		t.Helper()
	}
}
