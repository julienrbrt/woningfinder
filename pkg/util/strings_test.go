package util

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPadLeft(t *testing.T) {
	a := assert.New(t)
	a.Equal("BBBBBA", PadLeft("A", "B", 6))
	a.Equal("DDABC", PadLeft("ABC", "D", 5))
	a.Equal("ABCD", PadLeft("ABCD", "E", 4))
	a.Equal("ABCD", PadLeft("ABCD", "E", 3))
}

func TestPadRight(t *testing.T) {
	a := assert.New(t)
	a.Equal("ABBBBB", PadRight("A", "B", 6))
	a.Equal("ABCDD", PadRight("ABC", "D", 5))
	a.Equal("ABCD", PadRight("ABCD", "E", 4))
	a.Equal("ABCD", PadRight("ABCD", "E", 3))
}

func TestCamelCase(t *testing.T) {
	a := assert.New(t)
	a.Equal("ThisIsATest", CamelCase("this_is_a_test"))
	a.Equal("ThisIsATest", CamelCase("this_  Is_ __aTest"))
}

func TestSnakeCase(t *testing.T) {
	a := assert.New(t)
	a.Equal("this_is_a_test123", SnakeCase("thisIs-a test123"))
	a.Equal("this_is_12a_test", SnakeCase("  this_  Is_ __12aTest"))
}

func TestIsNumeric(t *testing.T) {
	a := assert.New(t)
	a.True(IsNumeric("1.0"))
	a.True(IsNumeric("1"))
	a.True(IsNumeric("1.0001"))
	a.True(IsNumeric("99999.0001"))
	a.True(IsNumeric("009.0001"))
	a.True(IsNumeric("1e6"))
	a.True(IsNumeric("1.0e6"))
	a.False(IsNumeric("1.0e1.5"))
	a.True(IsNumeric(".06"))
	a.True(IsNumeric("-1"))
	a.True(IsNumeric("+1"))
	a.True(IsNumeric("+1e9"))
	a.True(IsNumeric("+1e-9"))
	a.False(IsNumeric("a"))
	a.False(IsNumeric("0.1.0"))
	a.False(IsNumeric(""))
	a.False(IsNumeric(" 1"))
	a.False(IsNumeric("1 "))
	a.False(IsNumeric("e6"))
	a.False(IsNumeric("4e3e2"))
}

func BenchmarkIsNumeric(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsNumeric("a 240")
		IsNumeric("abc123")
		IsNumeric("0.4.")
		IsNumeric("0.4")
	}
}

func BenchmarkParseFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.ParseFloat("a 240", 64)
		strconv.ParseFloat("abc123", 64)
		strconv.ParseFloat("0.4.", 64)
		strconv.ParseFloat("0.4.", 64)
	}
}
