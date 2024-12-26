package utils

import (
	"unicode"
)

var (
	ranger = []*unicode.RangeTable{unicode.Latin, unicode.Cyrillic, unicode.ASCII_Hex_Digit, unicode.Punct, unicode.White_Space}
)

// IsStrPrintable checks if string can be printed to console withour artifacts
func IsStrPrintable(s string) bool {
	if len(s) > 1024 {
		s = s[:1024]
	}
	for _, r := range []rune(s) {
		if !unicode.IsOneOf(ranger, r) {
			return false
		}
	}
	return true
}
