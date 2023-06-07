package test

import (
	"encoding/json"
	store "github.com/launchquickly/ges/pkg"
	"github.com/stretchr/testify/require"
	"testing"
)

func AssertStreamContents(t *testing.T, fixture *Fixture, stream store.Stream) {
	t.Helper()

	for i := 0; i < fixture.ExpectedLength; i++ {
		AssertVersionsMatch(t, fixture.Records[offsetIndex(i, int(fixture.FromVersion))], stream[i])
		AssertJSONDataMatch(t, fixture.Records[offsetIndex(i, int(fixture.FromVersion))], stream[i])
	}
}

func AssertJSONDataMatch(t *testing.T, expected, actual store.Record) {
	t.Helper()

	e, a, err := marshalMsgsToStrings(expected.Data, actual.Data)
	require.NoError(t, err)

	require.JSONEq(t, e, a)
}

func AssertStreamLength(t *testing.T, stream store.Stream, length int) {
	t.Helper()

	require.NotNil(t, stream)
	require.Len(t, stream, length)
}

func AssertVersionsMatch(t *testing.T, expected, actual store.Record) {
	t.Helper()

	require.Equal(t, expected.Version, actual.Version, "version mismatch")
}

func marshalMsgsToStrings(msg1, msg2 json.RawMessage) (string, string, error) {
	m1, err := json.Marshal(msg1)
	if err != nil {
		return "", "", err
	}
	m2, err := json.Marshal(msg2)
	if err != nil {
		return "", "", err
	}
	return string(m1), string(m2), nil
}

func offsetIndex(i, from int) int {
	if from == 0 {
		return i
	}
	return (from - 1) + i
}
