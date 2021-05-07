package util

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"strings"
	"unicode"
)

func PadRight(str, pad string, length int) string {
	for len(str) < length {
		str += pad
	}
	return str
}

func PadLeft(str, pad string, length int) string {
	for len(str) < length {
		str = pad + str
	}
	return str
}

func CamelCase(input string) string {
	inputWithSpaces := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(input, "_", " "), "-", " "))
	result := make([]rune, 0, len(inputWithSpaces))
	reader := bytes.NewReader([]byte(inputWithSpaces))
	first := true
	for {
		r, err := readNextCamelCaseRune(reader, first)
		if err != nil {
			break
		}
		result = append(result, r)
		first = false
	}

	return string(result)
}

func readNextCamelCaseRune(reader *bytes.Reader, shouldUpperCase bool) (rune, error) {
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return 0, err
		}
		if unicode.IsSpace(r) {
			shouldUpperCase = true
			continue
		}
		if shouldUpperCase {
			r = unicode.ToUpper(r)
		}
		return r, nil
	}
}

func SnakeCase(input string) string {
	result := make([]rune, 0, len(input))
	reader := bytes.NewReader([]byte(input))
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		if unicode.IsUpper(r) {
			result = startNewWord(result)
			r = unicode.ToLower(r)
		} else if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			result = startNewWord(result)
			continue
		}
		result = append(result, r)
	}

	return string(result)
}

func startNewWord(result []rune) []rune {
	if len(result) > 0 && result[len(result)-1] != '_' {
		result = append(result, '_')
	}
	return result
}

// IsNumeric checks if a given string can be considered numeric
// The following formats will be considered numeric:
// * All Integers
// * Floating Point e.g. 21.4
// * e notation e.g. 1e6
// * All of the above when prefixed with either + or -
func IsNumeric(s string) bool {
	return isNumeric(s, false)
}

func isNumeric(s string, isExponent bool) bool {
	start := 0
	if len(s) == 0 {
		return false
	}
	hasDot := false
	switch s[0] {
	case 'e':
		return false
	case '-', '+':
		start = 1
	}
	for i, r := range s[start:] {
		switch r {
		case '.':
			if isExponent || hasDot {
				return false
			}
			hasDot = true
		case 'e':
			if isExponent {
				return false
			}
			return isNumeric(s[i+start+1:], true)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			continue
		default:
			return false
		}
	}
	return true
}

func RandomString(length int) (str string) {
	b := make([]byte, length)
	rand.Read(b)
	str = fmt.Sprintf("%x", b)[:length]
	return
}
