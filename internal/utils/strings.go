package utils

import (
	"unicode"
)

var (
	ranger = []*unicode.RangeTable{
		unicode.Latin,
		unicode.Cyrillic,
		unicode.Punct,
		unicode.Digit,
		unicode.White_Space,
		unicode.Quotation_Mark,
		unicode.Hyphen,
		unicode.Pattern_Syntax,
		unicode.Symbol,
	}
)

const (
	PrintableStrLen = 2048
)

// IsStrPrintable checks if string can be printed to console withour artifacts
func IsStrPrintable(s string) bool {
	if len(s) > PrintableStrLen {
		s = s[:PrintableStrLen]
	}
	for _, r := range []rune(s) {
		if !unicode.IsOneOf(ranger, r) {
			return false
		}
	}
	return true
}
