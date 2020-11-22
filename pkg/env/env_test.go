package env

import (
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_GetString(t *testing.T) {
	os.Setenv("foo", "bar")

	v, err := GetString("foo")
	if err != nil {
		t.Errorf("Expected GetString to return value and no error, got %s", err)
	}
	if v != "bar" {
		t.Errorf("Expected GetString to return 'bar', got %s", v)
	}

	if _, err := GetString("unknown"); err == nil {
		t.Error("Expected GetString to return error, got none")
	}
}

func TestConfig_GetStringOrDefault(t *testing.T) {
	os.Setenv("foo", "bar")

	if v := GetStringOrDefault("foo", "default"); v != "bar" {
		t.Errorf("Expected GetStringOrDefault to return bar', got %s", v)
	}
	if v := GetStringOrDefault("unknown", "default"); v != "default" {
		t.Errorf("Expected GetStringOrDefault to return default', got %s", v)
	}
}

func TestConfig_MustGetStringSuccess(t *testing.T) {
	os.Setenv("foo", "bar")

	if v := MustGetString("foo"); v != "bar" {
		t.Errorf("Expected GetString to return 'bar', got %s", v)
	}
}

func TestConfig_MustGetStringFailure(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected MustGetString to panic, but it didn't")
		}
	}()

	os.Setenv("foo", "bar")
	MustGetString("unknown")

	t.Errorf("Expected MustGetString to panic, but it didn't")
}

func TestConfig_GetInt(t *testing.T) {
	os.Setenv("foo", "1")

	v, err := GetInt("foo")
	if err != nil {
		t.Errorf("Expected GetInt to return value and no error, got %s", err)
	}
	if v != 1 {
		t.Errorf("Expected GetInt to return '1', got %d", v)
	}

	if _, err := GetInt("unknown"); err == nil {
		t.Error("Expected GetInt to return error, got none")
	}
}

func TestConfig_GetIntOrDefault(t *testing.T) {
	os.Setenv("foo", "1")

	if v := GetIntOrDefault("foo", 2); v != 1 {
		t.Errorf("Expected GetIntOrDefault to return 1', got %d", v)
	}
	if v := GetIntOrDefault("unknown", 2); v != 2 {
		t.Errorf("Expected GetIntOrDefault to return 2', got %d", v)
	}
}

func TestConfig_MustGetIntSuccess(t *testing.T) {
	os.Setenv("foo", "1")

	if v := MustGetInt("foo"); v != 1 {
		t.Errorf("Expected GetInt to return '1', got %d", v)
	}
}

func TestConfig_MustGetIntFailure(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected MustGetInt to panic, but it didn't")
		}
	}()

	os.Setenv("foo", "1")
	MustGetInt("unknown")

	t.Errorf("Expected MustGetInt to panic, but it didn't")
}

func TestConfig_GetBool(t *testing.T) {
	env := map[string]string{
		"foo":     "true",
		"bar":     "F",
		"baz":     "1",
		"foobar":  "0",
		"invalid": "not a boolean value",
	}
	for k, v := range env {
		os.Setenv(k, v)
	}

	v, err := GetBool("foo")
	if err != nil {
		t.Errorf("Expected GetBool to return value and no error, got %s", err)
	}
	if !v {
		t.Errorf("Expected GetBool to return true', got %t", v)
	}

	v, err = GetBool("bar")
	if err != nil {
		t.Errorf("Expected GetBool to return value and no error, got %s", err)
	}
	if v {
		t.Errorf("Expected GetBool to return false', got %t", v)
	}

	v, err = GetBool("baz")
	if err != nil {
		t.Errorf("Expected GetBool to return value and no error, got %s", err)
	}
	if !v {
		t.Errorf("Expected GetBool to return true', got %t", v)
	}

	v, err = GetBool("foobar")
	if err != nil {
		t.Errorf("Expected GetBool to return value and no error, got %s", err)
	}
	if v {
		t.Errorf("Expected GetBool to return false', got %t", v)
	}

	if _, err := GetBool("invalid"); err == nil {
		t.Error("Expected GetBool to return error, got none")
	}

	if _, err := GetBool("unknown"); err == nil {
		t.Error("Expected GetBool to return error, got none")
	}
}

func TestConfig_GetBoolOrDefault(t *testing.T) {
	env := map[string]string{
		"foo":     "false",
		"bar":     "T",
		"invalid": "not a boolean value",
	}
	for k, v := range env {
		os.Setenv(k, v)
	}

	if v := GetBoolOrDefault("foo", true); v {
		t.Errorf("Expected GetBool to return false', got %t", v)
	}
	if v := GetBoolOrDefault("bar", false); !v {
		t.Errorf("Expected GetBool to return true', got %t", v)
	}
	if v := GetBoolOrDefault("invalid", true); !v {
		t.Errorf("Expected GetBool to return true', got %t", v)
	}
	if v := GetBoolOrDefault("unknown", false); v {
		t.Errorf("Expected GetBool to return false', got %t", v)
	}
}

func TestConfig_GetURL(t *testing.T) {
	env := map[string]string{
		"foo":     "https://www.sqills.com/path",
		"invalid": "",
	}
	for k, v := range env {
		os.Setenv(k, v)
	}

	expectedFoo := &url.URL{
		Scheme: "https",
		Host:   "www.sqills.com",
		Path:   "/path",
	}

	v, err := GetURL("foo")
	if err != nil {
		t.Errorf("Expected GetURL to return value and no error, got %s", err)
	}
	if v.String() != expectedFoo.String() {
		t.Errorf("Expected GetURL to return '%s', got %s", expectedFoo, v)
	}

	if v, err := GetURL("invalid"); err == nil {
		t.Errorf("Expected GetURL to return error, got none with value %s", v)
	}

	if v, err := GetURL("unknown"); err == nil {
		t.Errorf("Expected GetURL to return error, got none with value %s", v)
	}
}

func TestConfig_MustGetURLSuccess(t *testing.T) {
	os.Setenv("foo", "https://www.sqills.com/path")

	expectedFoo := &url.URL{
		Scheme: "https",
		Host:   "www.sqills.com",
		Path:   "/path",
	}

	if v := MustGetURL("foo"); v.String() != expectedFoo.String() {
		t.Errorf("Expected GetURL to return '%s', got %s", expectedFoo, v)
	}
}

func TestConfig_MustGetStringInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected MustGetURL to panic, but it didn't")
		}
	}()

	os.Setenv("invalid", "")
	MustGetURL("invalid")

	t.Errorf("Expected MustGetURL to panic, but it didn't")
}

func TestConfig_MustGetURLUnknown(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected MustGetURL to panic, but it didn't")
		}
	}()

	MustGetURL("unknown")
	t.Errorf("Expected MustGetURL to panic, but it didn't")
}

func TestConfig_GetJSON(t *testing.T) {
	env := map[string]string{
		"valid":   `"test"`,
		"invalid": "{]",
	}
	for k, v := range env {
		os.Setenv(k, v)
	}

	a := assert.New(t)
	var v string
	a.Error(GetJSON("invalid", &v))
	a.Error(GetJSON("unknown", &v))
	a.NoError(GetJSON("valid", &v))
}

func TestConfig_MustGetJSON_Unknown(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected MustGetJSON to panic, but it didn't")
		}
	}()

	MustGetJSON("unknown", make(map[string]string))
	t.Errorf("Expected MustGetJSON to panic, but it didn't")
}

func TestConfig_MustGetJSON_InvalidJSON(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected MustGetJSON to panic, but it didn't")
		}
	}()

	os.Setenv("key", `{"invalid": "json"`)
	MustGetJSON("key", make(map[string]string))

	t.Errorf("Expected MustGetJSON to panic, but it didn't")
}

func TestConfig_MustGetJSON_Success(t *testing.T) {
	os.Setenv("key", `{"valid": "json", "foo": "bar"}`)

	var v map[string]string
	MustGetJSON("key", &v)
	a := assert.New(t)
	a.Equal("json", v["valid"])
	a.Equal("bar", v["foo"])
}

func TestConfig_GetStringList(t *testing.T) {
	os.Setenv("foo", "bar,baz,qux")

	v, err := GetStringList("foo")
	if err != nil {
		t.Errorf("Expected GetStringList to return value and no error, got %v", err)
	}
	if len(v) != 3 {
		t.Errorf("Expected GetStringList to return 3 strings, got %d", len(v))
	}

	if _, err := GetStringList("unknown"); err == nil {
		t.Error("Expected GetStringList to return error, got none")
	}
}

func TestConfig_GetStringListOrDefault(t *testing.T) {
	os.Setenv("foo", "bar,baz,qux")

	if v := GetStringListOrDefault("foo", []string{}); len(v) != 3 {
		t.Errorf("Expected GetStringListOrDefault to return 3 strings, got %d", len(v))
	}
	if v := GetStringListOrDefault("unknown", []string{"default"}); len(v) != 1 {
		t.Errorf("Expected GetStringListOrDefault to return 1 string, got %d", len(v))
	}
}

func TestConfig_MustGetStringListSuccess(t *testing.T) {
	os.Setenv("foo", "bar,baz")

	if v := MustGetStringList("foo"); len(v) != 2 {
		t.Errorf("Expected MustGetStringList to return 2 strings, got %d", len(v))
	}
}

func TestConfig_MustGetStringListFailure(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected MustGetStringList to panic, but it didn't")
		}
	}()

	os.Setenv("foo", "bar,baz")
	MustGetString("unknown")

	t.Errorf("Expected MustGetStringList to panic, but it didn't")
}
