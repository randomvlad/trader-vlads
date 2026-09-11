package stringutil

import "unicode"

func Truncate(value string, maxLength int) string {
	runes := []rune(value)
	if len(runes) <= maxLength {
		return value
	}
	return string(runes[:maxLength])
}

func FirstCharLowercase(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	return string(unicode.ToLower(runes[0]))
}
