package match

import (
	"sort"
)

type matchesByIJ []*Match

func (s matchesByIJ) Len() int {
	return len(s)
}

func (s matchesByIJ) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s matchesByIJ) Less(i, j int) bool {
	switch {
	case s[i].I < s[j].I:
		return true
	case s[i].I == s[j].I:
		return s[i].J < s[j].J
	default:
		return false
	}
}

func Sort(matches []*Match) {
	sort.Stable(matchesByIJ(matches))
}
