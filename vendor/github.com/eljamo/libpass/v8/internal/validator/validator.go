package validator

import (
	"slices"
	"unicode/utf8"
)

func HasElementWithLengthGreaterThanOne(s []string) bool {
	for _, str := range s {
		if utf8.RuneCountInString(str) > 1 {
			return true
		}
	}

	return false
}

func IsElementInSlice[T comparable](s []T, e T) bool {
	return slices.Contains(s, e)
}
