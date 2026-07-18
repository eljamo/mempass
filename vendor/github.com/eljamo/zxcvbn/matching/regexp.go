package matching

import (
	"regexp"
	"strconv"

	"github.com/eljamo/zxcvbn/match"
	"github.com/eljamo/zxcvbn/scoring"
)

const recentYearPastWindow = 150

type regexpMatch struct {
	regexes []struct {
		Name   string
		Regexp *regexp.Regexp
	}
}

func (r regexpMatch) Matches(password string) []*match.Match {
	var matches []*match.Match
	for _, rx := range r.regexes {
		for _, indexes := range rx.Regexp.FindAllStringIndex(password, -1) {
			token := password[indexes[0]:indexes[1]]
			if rx.Name == "recent_year" && !isRecentYear(token) {
				continue
			}
			matches = append(matches, &match.Match{
				Pattern:   "regex",
				Token:     token,
				I:         indexes[0],
				J:         indexes[1] - 1,
				RegexName: rx.Name,
			})
		}
	}
	match.Sort(matches)
	return matches
}

func isRecentYear(token string) bool {
	year, err := strconv.Atoi(token)
	if err != nil {
		return false
	}
	return year >= scoring.ReferenceYear-recentYearPastWindow &&
		year <= scoring.ReferenceYear+scoring.MinYearSpace
}
