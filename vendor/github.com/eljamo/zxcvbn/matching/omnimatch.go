package matching

import (
	"regexp"

	"github.com/eljamo/zxcvbn/adjacency"
	"github.com/eljamo/zxcvbn/frequency"
	"github.com/eljamo/zxcvbn/match"
)

type Omnimatcher struct {
	dictMatcher dictionaryMatch
	l33tMatcher l33tMatch
}

func NewOmnimatcher(customDictionaries map[string][]string) Omnimatcher {
	dictMatcher := defaultRankedDictionaries
	for name, words := range customDictionaries {
		dictMatcher = dictMatcher.withDict(name, buildRankedDict(words))
	}

	return Omnimatcher{
		dictMatcher: dictMatcher,
		l33tMatcher: newl33tMatch(dictMatcher, l33tTable),
	}
}

func Omnimatch(password string, userInputs []string) (matches []*match.Match) {
	return defaultOmnimatcher.Omnimatch(password, userInputs)
}

func (om Omnimatcher) Omnimatch(password string, userInputs []string) (matches []*match.Match) {
	dictMatcher := om.dictMatcher
	l33tMatchers := []match.Matcher{om.l33tMatcher}

	if len(userInputs) > 0 {
		userInputDict := buildRankedDict(userInputs)
		dictMatcher = om.dictMatcher.withDict("user_inputs", userInputDict)
		userInputMatcher := dictionaryMatch{
			rankedDictionaries: map[string]rankedDictionnary{
				"user_inputs": userInputDict,
			},
		}
		l33tMatchers = append(l33tMatchers, newl33tMatch(userInputMatcher, l33tTable))
	}

	matchers := []match.Matcher{
		dictMatcher,
		reverseDictionnaryMatch{dm: dictMatcher},
	}

	matchers = append(matchers, l33tMatchers...)

	matchers = append(matchers,
		spatialMatch{graphs: defaultGraphs},
		repeatMatch{omnimatch: om.Omnimatch, userInputs: userInputs},
		sequenceMatch{},
		regexpMatch{regexes: defaultRegexpMatch},
		dateMatch{},
	)

	for _, m := range matchers {
		matches = append(matches, m.Matches(password)...)
	}
	match.Sort(matches)
	return matches
}

var (
	defaultRankedDictionaries = loadDefaultDictionaries()
	defaultGraphs             = loadDefaultAdjacencyGraphs()
	defaultRegexpMatch        = []struct {
		Name   string
		Regexp *regexp.Regexp
	}{
		{
			Name:   "recent_year",
			Regexp: regexp.MustCompile(`\d{4}`),
		},
	}
	l33tTable = map[string][]string{
		"a": {"4", "@"},
		"b": {"8"},
		"c": {"(", "{", "[", "<"},
		"e": {"3"},
		"g": {"6", "9"},
		"i": {"1", "!", "|"},
		"l": {"1", "|", "7"},
		"o": {"0"},
		"s": {"$", "5"},
		"t": {"+", "7"},
		"x": {"%"},
		"z": {"2"},
	}
)

var defaultOmnimatcher = NewOmnimatcher(nil)

func loadDefaultDictionaries() dictionaryMatch {
	rd := make(map[string]rankedDictionnary)
	for n, list := range frequency.FrequencyLists {
		rd[n] = buildRankedDict(list)
	}
	return dictionaryMatch{
		rankedDictionaries: rd,
	}
}

func loadDefaultAdjacencyGraphs() []*adjacency.Graph {
	return []*adjacency.Graph{
		adjacency.Graphs["qwerty"],
		adjacency.Graphs["dvorak"],
		adjacency.Graphs["keypad"],
		adjacency.Graphs["mac_keypad"],
	}
}
