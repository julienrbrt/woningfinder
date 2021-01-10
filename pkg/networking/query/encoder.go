package query

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/woningfinder/woningfinder/pkg/util"
)

const Tag = "query"

type QueryEncoder interface {
	Query() (Query, error)
}

type ValueEncoder interface {
	QueryValue() (string, error)
}

type Encoder interface {
	Encode(interface{}) (Query, error)
}

type encoder struct{}

// NewEncoder creates an Encoder
func NewEncoder() Encoder {
	return encoder{}
}

func (e encoder) Encode(i interface{}) (Query, error) {
	return e.encodeQuery(reflect.ValueOf(i))
}

func (e encoder) encodeQuery(value reflect.Value) (Query, error) {
	if !value.IsValid() || value.IsZero() {
		return make(Query, 0), nil
	}
	if v, ok := value.Interface().(QueryEncoder); ok {
		return v.Query()
	}

	switch value.Kind() {
	case reflect.Ptr, reflect.Interface:
		return e.encodeQuery(value.Elem())
	case reflect.Struct:
		return e.encodeStruct(value)
	case reflect.Map:
		return e.encodeMap(value)
	default:
		return nil, fmt.Errorf("unable to encode query for type %q", value.Kind())
	}
}

func (e encoder) encodeStruct(value reflect.Value) (Query, error) {
	q := make(Query, 0)
	for i := 0; i < value.NumField(); i++ {
		if !value.Field(i).CanSet() {
			continue
		}
		key := e.keyForField(value.Type().Field(i))
		if key == "-" {
			continue
		}
		if err := e.encodeValue(&q, key, value.Field(i)); err != nil {
			return nil, err
		}
	}

	return q, nil
}

func (e encoder) encodeMap(value reflect.Value) (Query, error) {
	q := make(Query, 0)
	for _, key := range value.MapKeys() {
		if err := e.encodeValue(&q, e.keyForMapKey(key), value.MapIndex(key)); err != nil {
			return nil, err
		}
	}

	return q, nil
}

func (e encoder) encodeValue(q *Query, key string, value reflect.Value) error {
	if value.IsZero() {
		return nil
	}
	if v, ok := value.Interface().(ValueEncoder); ok {
		value, err := v.QueryValue()
		if err != nil {
			return err
		}
		q.Add(key, value)
		return nil
	}

	if str, ok := e.scalarValueToString(value); ok {
		q.Add(key, str)
		return nil
	}

	switch value.Kind() {
	case reflect.Ptr, reflect.Interface:
		return e.encodeValue(q, key, value.Elem())
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			if err := e.encodeValue(q, key, value.Index(i)); err != nil {
				return err
			}
		}
	default:
		q2, err := e.encodeQuery(value)
		if err != nil {
			return err
		}
		log.Printf("Encoded query for key %s: %s", key, q2.Encode())
		e.merge(q, key, q2)
	}
	return nil
}

func (e encoder) merge(q *Query, key string, query Query) {
	for _, p := range query {
		keyParts := strings.SplitN(p.Key, "[", 2)
		newKey := fmt.Sprintf("%s[%s]", key, keyParts[0])
		if len(keyParts) > 1 {
			newKey = fmt.Sprintf("%s[%s", newKey, keyParts[1])
		}
		q.Add(newKey, p.Value)
	}
}

func (e encoder) keyForField(field reflect.StructField) string {
	key := field.Tag.Get(Tag)
	if key == "" {
		key = util.SnakeCase(field.Name)
	}
	return key
}

func (e encoder) keyForMapKey(key reflect.Value) string {
	if s, ok := key.Interface().(fmt.Stringer); ok {
		return s.String()
	}
	if str, ok := e.scalarValueToString(key); ok {
		return str
	}
	return ""
}

func (e encoder) scalarValueToString(value reflect.Value) (string, bool) {
	switch value.Kind() {
	case reflect.String:
		return value.String(), true
	case reflect.Bool:
		return strconv.FormatBool(value.Bool()), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10), true
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'G', -1, 64), true
	}
	return "", false
}
