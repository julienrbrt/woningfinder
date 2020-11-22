package env

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// GetString gets an environment variable value as a string
func GetString(variable string) (string, error) {
	value := os.Getenv(variable)
	if value != "" {
		return value, nil
	}

	return "", fmt.Errorf("unable to retrieve variable %q", variable)
}

// MustGetString gets an environment variable value as an string
// It panics if a it does not exists
func MustGetString(variable string) string {
	val, err := GetString(variable)
	if err != nil {
		panic(fmt.Errorf("failed to get environment variable %q as string: %w", variable, err))
	}
	return val
}

// GetStringOrDefault gets an environment variable value as a string
// If unexisting default to the given value
func GetStringOrDefault(variable, def string) string {
	if val, err := GetString(variable); err == nil {
		return val
	}
	return def
}

// GetInt gets an environment variable value as an int
func GetInt(variable string) (int, error) {
	if value, err := GetString(variable); err == nil && value != "" {
		return strconv.Atoi(value)
	}

	return 0, fmt.Errorf("unable to retrieve variable %q", variable)
}

// MustGetInt gets an environment variable value as an int
// It panics if a it does not exists
func MustGetInt(variable string) int {
	val, err := GetInt(variable)
	if err != nil {
		panic(fmt.Errorf("failed to get environment variable %q as int: %w", variable, err))
	}
	return val
}

// GetIntOrDefault gets an environment variable value as a int
// If unexisting default to the given value
func GetIntOrDefault(variable string, def int) int {
	if val, err := GetInt(variable); err == nil {
		return val
	}
	return def
}

// GetBool gets an environment variable value as a bool
func GetBool(variable string) (bool, error) {
	val, err := GetString(variable)
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(val)
}

// GetBoolOrDefault gets an environment variable value as a bool
// If unexisting default to the given value
func GetBoolOrDefault(variable string, def bool) bool {
	if val, err := GetBool(variable); err == nil {
		return val
	}
	return def
}

// GetURL gets an environment variable value as an url.URL
func GetURL(variable string) (*url.URL, error) {
	val, err := GetString(variable)
	if err != nil {
		return nil, err
	}

	return url.Parse(val)
}

// MustGetURL gets an environment variable value as an url.URL
// It panics if a it does not exists
func MustGetURL(variable string) *url.URL {
	val, err := GetURL(variable)
	if err != nil {
		panic(fmt.Errorf("failed to get environment variable %q as url: %w", variable, err))
	}
	return val
}

// GetJSON gets an environment variable value as a JSON
func GetJSON(variable string, value interface{}) error {
	val, err := GetString(variable)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(val), value); err != nil {
		return fmt.Errorf("failed to unmarshal variable %q from JSON: %w", variable, err)
	}

	return nil
}

// MustGetJSON gets an environment variable value as a JSON
// It panics if a it does not exists
func MustGetJSON(variable string, value interface{}) {
	if err := GetJSON(variable, value); err != nil {
		panic(fmt.Errorf("failed to get environment variable %q as JSON: %w", variable, err))
	}
}

// GetStringList gets an environment variable value as an array of string
func GetStringList(variable string) ([]string, error) {
	str, err := GetString(variable)
	if err != nil {
		return nil, err
	}
	return strings.Split(str, ","), nil
}

// MustGetStringList gets an environment variable value as an array of string
// It panics if a it does not exists
func MustGetStringList(variable string) []string {
	return strings.Split(MustGetString(variable), ",")
}

// GetStringListOrDefault gets an environment variable value as an array of string
// If unexisting default to the given value
func GetStringListOrDefault(variable string, def []string) []string {
	if v := GetStringOrDefault(variable, ""); v != "" {
		return strings.Split(v, ",")
	}
	return def
}
