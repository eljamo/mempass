package matching

import (
	"strings"
	"unicode"

	"github.com/eljamo/zxcvbn/match"
)

type dictionaryMatch struct {
	rankedDictionaries map[string]rankedDictionnary
}

// lowerAligned lowercases s once and returns an offset table mapping each
// rune-start byte index of s (plus len(s)) to the corresponding byte index in
// the lowered string. Lowering can change a rune's UTF-8 length, so original
// indices cannot be used in the lowered string directly.
func lowerAligned(s string) (string, []int) {
	var b strings.Builder
	b.Grow(len(s))
	offsets := make([]int, len(s)+1)
	for i, r := range s {
		offsets[i] = b.Len()
		b.WriteRune(unicode.ToLower(r))
	}
	offsets[len(s)] = b.Len()
	return b.String(), offsets
}

func (dm dictionaryMatch) Matches(password string) []*match.Match {
	var results []*match.Match

	// lowercase once; per-substring ToLower makes this loop O(n^3)
	lower, lowerOffset := lowerAligned(password)

	for dictionaryName, rankedDict := range dm.rankedDictionaries {
		for i := range password {
			j := len(password) - 1
			for delta := range password[i:] {
				if delta > 0 {
					j = i + delta - 1
				}
				word := lower[lowerOffset[i]:lowerOffset[j+1]]
				if val, ok := rankedDict[word]; ok {
					matchDic := &match.Match{
						Pattern:        "dictionary",
						I:              i,
						J:              j,
						Token:          password[i : j+1],
						MatchedWord:    word,
						Rank:           val,
						DictionaryName: dictionaryName,
					}

					results = append(results, matchDic)
				}
			}
		}
	}

	match.Sort(results)
	return results
}

func (dm dictionaryMatch) withDict(name string, d rankedDictionnary) dictionaryMatch {
	rd2 := make(map[string]rankedDictionnary, len(dm.rankedDictionaries)+1)
	for k, v := range dm.rankedDictionaries {
		rd2[k] = v
	}
	rd2[name] = d
	return dictionaryMatch{rankedDictionaries: rd2}
}

type rankedDictionnary map[string]int

func buildRankedDict(unrankedList []string) rankedDictionnary {
	result := make(rankedDictionnary)

	for i, v := range unrankedList {
		result[strings.ToLower(v)] = i + 1
	}

	return result
}
