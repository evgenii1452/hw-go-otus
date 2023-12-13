package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var result strings.Builder
	var prevRune rune

	if err := validate(str); err != nil {
		return "", err
	}

	for i, currentRune := range str {
		if i == 0 {
			prevRune = currentRune
			continue
		}

		if value, _ := strconv.Atoi(string(currentRune)); unicode.IsDigit(currentRune) && value > 0 {
			result.WriteString(strings.Repeat(string(prevRune), value))
		}

		if !unicode.IsDigit(currentRune) && !unicode.IsDigit(prevRune) {
			result.WriteString(string(prevRune))
		}

		if i == len(str)-1 && !unicode.IsDigit(currentRune) {
			result.WriteString(string(currentRune))
		}

		prevRune = currentRune
	}

	return result.String(), nil
}

func validate(str string) error {
	var prevRune rune

	if len(str) > 0 && unicode.IsDigit(rune(str[0])) {
		return ErrInvalidString
	}

	for i, currentRune := range str {
		if i != 0 && unicode.IsDigit(currentRune) && unicode.IsDigit(prevRune) {
			return ErrInvalidString
		}

		prevRune = currentRune
	}

	return nil
}
