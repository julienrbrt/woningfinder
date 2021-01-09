package query_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/networking/query"
	"github.com/woningfinder/woningfinder/internal/util"
)

type Foo struct {
	Bar           string `query:"bar_key"`
	BarPtr        *string
	Baz           int
	Foobar        *Foobar
	Foobaz        int `query:"-"`
	ignore        string
	ScalarBool    bool
	ScalarFloat64 float64
	ScalarUint    uint32
	Map           map[StringerMapKey]float32
	CustomQuery   CustomQuery
	CustomMap     map[string]CustomQueryValue
	CustomSlice   []CustomQueryValue
}

type StringerMapKey struct{}

func (k StringerMapKey) String() string {
	return "map_key"
}

type Foobar struct {
	Foob             []string         `query:"fo_ob"`
	Nested           Nested           `query:"nested"`
	Map              map[int]string   `query:"map"`
	CustomQueryValue CustomQueryValue `query:"custom_query_value"`
}

type Nested struct {
	Foo string `query:"foobar"`
}

type CustomQueryValue struct {
	err   error
	value string
}

func (v CustomQueryValue) QueryValue() (string, error) {
	return v.value, v.err
}

type CustomQuery struct {
	err   error
	query query.Query
}

func (v CustomQuery) Query() (query.Query, error) {
	return v.query, v.err
}

func TestEncoder_Encode_CustomQueryValueFailure(t *testing.T) {
	a := assert.New(t)
	data := Foo{
		Foobar: &Foobar{
			CustomQueryValue: CustomQueryValue{
				err: fmt.Errorf("error while getting custom query value"),
			},
		},
	}
	_, err := query.NewEncoder().Encode(&data)
	a.Error(err)
	a.True(errors.Is(err, data.Foobar.CustomQueryValue.err))
}

func TestEncoder_Encode_CustomQueryFailure(t *testing.T) {
	a := assert.New(t)
	data := Foo{
		CustomQuery: CustomQuery{
			err: fmt.Errorf("error while getting custom query"),
		},
	}
	_, err := query.NewEncoder().Encode(&data)
	a.Error(err)
	a.True(errors.Is(err, data.CustomQuery.err))
}

func TestEncoder_Encode_CustomMapFailure(t *testing.T) {
	a := assert.New(t)
	data := Foo{
		CustomMap: map[string]CustomQueryValue{
			"err": {err: fmt.Errorf("error while getting custom query value")},
		},
	}
	_, err := query.NewEncoder().Encode(&data)
	a.Error(err)
	a.True(errors.Is(err, data.CustomMap["err"].err))
}

func TestEncoder_Encode_CustomSliceFailure(t *testing.T) {
	a := assert.New(t)
	data := Foo{
		CustomSlice: []CustomQueryValue{
			{err: fmt.Errorf("error while getting custom query value")},
		},
	}
	_, err := query.NewEncoder().Encode(&data)
	a.Error(err)
	a.True(errors.Is(err, data.CustomSlice[0].err))
}

func TestEncoder_Encode_Nil(t *testing.T) {
	a := assert.New(t)
	q, err := query.NewEncoder().Encode(nil)
	a.NoError(err)
	a.Len(q, 0)
}

func TestEncoder_Encode_NonStruct(t *testing.T) {
	a := assert.New(t)
	_, err := query.NewEncoder().Encode("non-struct")
	a.Error(err)
}

func TestEncoder_Encode_Success(t *testing.T) {
	a := assert.New(t)
	data := Foo{
		Bar:    "bar-value",
		BarPtr: util.StrPtr("bar-ptr-value"),
		Baz:    4,
		Foobar: &Foobar{
			Foob: []string{"test1", "test2", "test3"},
			Nested: Nested{
				Foo: "bar",
			},
			Map: map[int]string{1: "item1", 2: "item2"},
			CustomQueryValue: CustomQueryValue{
				value: "custom_value",
			},
		},
		Foobaz:        5,
		ScalarBool:    true,
		ScalarFloat64: 14.2,
		ScalarUint:    124,
		Map: map[StringerMapKey]float32{
			StringerMapKey{}: 16.25,
		},
		CustomQuery: CustomQuery{
			query: query.Query{
				{Key: "custom_foo", Value: "custom_bar"},
				{Key: "custom_foo2", Value: "custom_bar2"},
			},
		},
		CustomMap: map[string]CustomQueryValue{
			"foo": {value: "bar"},
		},
		CustomSlice: []CustomQueryValue{
			{value: "foo"},
			{value: "bar"},
		},
	}
	q, err := query.NewEncoder().Encode(&data)
	a.NoError(err)
	a.Contains(q, query.Pair{Key: "bar_key", Value: "bar-value"})
	a.Contains(q, query.Pair{Key: "bar_ptr", Value: "bar-ptr-value"})
	a.Contains(q, query.Pair{Key: "baz", Value: "4"})
	a.Contains(q, query.Pair{Key: "foobar[fo_ob]", Value: "test1"})
	a.Contains(q, query.Pair{Key: "foobar[fo_ob]", Value: "test2"})
	a.Contains(q, query.Pair{Key: "foobar[fo_ob]", Value: "test3"})
	a.Contains(q, query.Pair{Key: "foobar[nested][foobar]", Value: "bar"})
	a.Contains(q, query.Pair{Key: "foobar[map][1]", Value: "item1"})
	a.Contains(q, query.Pair{Key: "foobar[map][2]", Value: "item2"})
	a.Contains(q, query.Pair{Key: "foobar[custom_query_value]", Value: "custom_value"})
	a.Contains(q, query.Pair{Key: "scalar_bool", Value: "true"})
	a.Contains(q, query.Pair{Key: "scalar_float64", Value: "14.2"})
	a.Contains(q, query.Pair{Key: "scalar_uint", Value: "124"})
	a.Contains(q, query.Pair{Key: "map[map_key]", Value: "16.25"})
	a.Contains(q, query.Pair{Key: "custom_query[custom_foo]", Value: "custom_bar"})
	a.Contains(q, query.Pair{Key: "custom_query[custom_foo2]", Value: "custom_bar2"})
	a.Contains(q, query.Pair{Key: "custom_map[foo]", Value: "bar"})
	a.Contains(q, query.Pair{Key: "custom_slice", Value: "foo"})
	a.Contains(q, query.Pair{Key: "custom_slice", Value: "bar"})
	a.NotContains(q, query.Pair{Key: "foobaz", Value: "5"})
}
