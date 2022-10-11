package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var tmpChar rune
	var hasTmpChar bool
	b := strings.Builder{}

	for _, char := range s {
		switch {
		case unicode.IsDigit(char):
			if !hasTmpChar {
				return "", ErrInvalidString
			}

			quantity, _ := strconv.Atoi(string(char))

			for i := 0; i < quantity; i++ {
				b.WriteRune(tmpChar)
			}
			hasTmpChar = false
		default:
			if hasTmpChar {
				b.WriteRune(tmpChar)
			}

			tmpChar = char
			hasTmpChar = true
		}
	}

	if hasTmpChar {
		b.WriteRune(tmpChar)
	}

	return b.String(), nil
}
