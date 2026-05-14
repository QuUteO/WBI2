package unpack

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// ErrInvalidString возвращается, когда строка имеет неверный формат
// (начинается с цифры или содержит только цифры)
var ErrInvalidString = errors.New("некорректная строка")

// Unpack распаковывает строку с повторяющимися символами.
// Поддерживает escape-последовательности через обратный слеш.
// Пример: "a4bc2d5e" -> "aaaabccddddde"
func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	var result strings.Builder
	runes := []rune(s)
	length := len(runes)

	for i := 0; i < length; i++ {
		current := runes[i]

		if current == '\\' {
			if i+1 >= length {
				return "", ErrInvalidString
			}

			escaped := runes[i+1]

			if i+2 < length && unicode.IsDigit(runes[i+2]) {
				j := i + 2
				numStr := ""
				for j < length && unicode.IsDigit(runes[j]) {
					numStr += string(runes[j])
					j++
				}

				count, err := strconv.Atoi(numStr)
				if err != nil {
					return "", ErrInvalidString
				}

				if count > 0 {
					result.WriteString(strings.Repeat(string(escaped), count))
				}
				i = j - 1
			} else {
				result.WriteRune(escaped)
				i++
			}
			continue
		}

		if i == 0 && unicode.IsDigit(current) {
			return "", ErrInvalidString
		}

		if !unicode.IsDigit(current) {
			if i+1 < length && unicode.IsDigit(runes[i+1]) {
				j := i + 1
				numStr := ""
				for j < length && unicode.IsDigit(runes[j]) {
					numStr += string(runes[j])
					j++
				}

				count, err := strconv.Atoi(numStr)
				if err != nil {
					return "", ErrInvalidString
				}

				if count > 0 {
					result.WriteString(strings.Repeat(string(current), count))
				}
				i = j - 1
			} else {
				result.WriteRune(current)
			}
		} else {
			return "", ErrInvalidString
		}
	}

	return result.String(), nil
}
