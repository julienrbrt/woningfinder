package query

import (
	"net/url"
	"strings"
)

// Query is list of QueryPair
type Query []Pair

// Pair contains the query field/name (key) and its value.
type Pair struct {
	Key      string
	Value    string
	RawKey   string
	RawValue string
}

// Add permits to append a QueryPair to the Query.
// The key and value will be query encoded during the encoding of the query
func (q *Query) Add(key, value string) {
	newQuery := Pair{
		Key:   key,
		Value: value,
	}
	*q = append(*q, newQuery)
}

// AddRaw permits to append a QueryPair to the Query
// The key and value will NOT be query encoded during the encoding of the query
func (q *Query) AddRaw(key, value string) {
	newQuery := Pair{
		RawKey:   key,
		RawValue: value,
	}
	*q = append(*q, newQuery)
}

// Encode performs URL encoding of the query
func (q *Query) Encode() string {
	if len(*q) == 0 {
		return ""
	}
	var buf strings.Builder

	for _, query := range *q {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		key := query.RawKey
		if query.Key != "" {
			key = url.QueryEscape(query.Key)
		}
		buf.WriteString(key)
		buf.WriteByte('=')
		value := query.RawValue
		if query.Value != "" {
			value = url.QueryEscape(query.Value)
		}
		buf.WriteString(value)
	}

	return buf.String()
}
