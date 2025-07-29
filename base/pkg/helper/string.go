package helper

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func Capitalize(str string) string {
	if str == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(str)
	return string(unicode.ToUpper(r)) + strings.ToLower(str[size:])
}