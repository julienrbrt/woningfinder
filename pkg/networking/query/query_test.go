package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery_Encode(t *testing.T) {
	a := assert.New(t)
	query := Query{
		Pair{
			Key:   "filter_field",
			Value: "foo",
		},
		Pair{
			Key:   "filter_value",
			Value: "bar",
		},
		Pair{
			Key:   "filter_field",
			Value: "baz",
		},
		Pair{
			Key:   "filter_value",
			Value: "quux",
		},
		Pair{
			RawKey:   "raw/key",
			RawValue: "raw/value",
		},
	}
	a.Equal("filter_field=foo&filter_value=bar&filter_field=baz&filter_value=quux&raw/key=raw/value", query.Encode())
}

func TestQuery_Encode_SpecialChars(t *testing.T) {
	a := assert.New(t)
	query := Query{

		Pair{
			Key:   "special_char_value",
			Value: "?!$*?üö",
		},
		Pair{
			Key:   "special_field_ :|",
			Value: "val",
		},
	}
	a.Equal("special_char_value=%3F%21%24%2A%3F%C3%BC%C3%B6&special_field_+%3A%7C=val", query.Encode())
}

func TestQuery_Add(t *testing.T) {
	a := assert.New(t)
	expectedQuery := Query{
		Pair{
			Key:   "key1",
			Value: "val1",
		},
		Pair{
			Key:   "key2",
			Value: "val2",
		},
		Pair{
			Key:   "filter_field",
			Value: "foo",
		},
		Pair{
			Key:   "filter_value",
			Value: "bar",
		},
		Pair{
			RawKey:   "raw/key",
			RawValue: "raw/value",
		},
	}

	q := Query{}
	q.Add("key1", "val1")
	q.Add("key2", "val2")
	q.Add("filter_field", "foo")
	q.Add("filter_value", "bar")
	q.AddRaw("raw/key", "raw/value")
	a.Equal(expectedQuery, q)
}

func TestQuery_Encode_EmptyQuery(t *testing.T) {
	a := assert.New(t)
	query := Query{}
	a.Equal("", query.Encode())
}
