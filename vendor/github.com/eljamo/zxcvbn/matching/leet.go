package matching

import (
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/eljamo/zxcvbn/match"
)

type l33tMatch struct {
	dm            dictionaryMatch
	table         map[string][]string
	trieRoot      *dictionaryTrieNode
	reversedTable map[rune][]rune
}

func newl33tMatch(dm dictionaryMatch, table map[string][]string) l33tMatch {
	return l33tMatch{
		dm:            dm,
		table:         table,
		trieRoot:      buildDictionaryTrie(dm.rankedDictionaries),
		reversedTable: reverseL33tTable(table),
	}
}

func (lm l33tMatch) Matches(password string) []*match.Match {
	matches := []*match.Match{}

	root := lm.trieRoot
	if root == nil {
		root = buildDictionaryTrie(lm.dm.rankedDictionaries)
	}
	reversed := lm.reversedTable
	if reversed == nil {
		reversed = reverseL33tTable(lm.table)
	}
	runes := []rune(password)

	// byteStart[k] is the byte offset of the k-th rune;
	// byteStart[len(runes)] == len(password).
	byteStart := make([]int, len(runes)+1)
	b := 0
	for k, r := range runes {
		byteStart[k] = b
		b += utf8.RuneLen(r)
	}
	byteStart[len(runes)] = b

	for start := range runes {
		states := []l33tState{{node: root}}

		for end := start; end < len(runes); end++ {
			chr := runes[end]
			lower := unicode.ToLower(chr)
			nextStates := []l33tState{}

			for _, state := range states {
				if child := state.node.children[lower]; child != nil {
					nextState := state.clone()
					nextState.node = child
					nextStates = append(nextStates, nextState)
				}

				for _, plain := range reversed[chr] {
					child := state.node.children[plain]
					if child == nil {
						continue
					}

					nextState, ok := state.withSubstitution(chr, plain)
					if !ok {
						continue
					}
					nextState.node = child
					nextStates = append(nextStates, nextState)
				}
			}

			if len(nextStates) == 0 {
				break
			}
			nextStates = dedupL33tStates(nextStates)

			for _, state := range nextStates {
				if len(state.sub) == 0 {
					continue
				}

				token := string(runes[start : end+1])
				if len(token) <= 1 {
					// filter single-character l33t matches to reduce noise.
					// otherwise '1' matches 'i', '4' matches 'a', both very common English words
					continue
				}

				for _, entry := range state.node.entries {
					matches = append(matches, &match.Match{
						Pattern:        "dictionary",
						Token:          token,
						MatchedWord:    entry.word,
						Rank:           entry.rank,
						DictionaryName: entry.dictionaryName,
						I:              byteStart[start],
						J:              byteStart[end+1] - 1,
						L33t:           true,
						Sub:            stringSubstitution(state.sub),
					})
				}
			}

			states = nextStates
		}
	}

	match.Sort(matches)
	return matches
}

type dictionaryTrieNode struct {
	id       int
	children map[rune]*dictionaryTrieNode
	entries  []dictionaryTrieEntry
}

type dictionaryTrieEntry struct {
	word           string
	rank           int
	dictionaryName string
}

func buildDictionaryTrie(dicts map[string]rankedDictionnary) *dictionaryTrieNode {
	nextID := 0
	newNode := func() *dictionaryTrieNode {
		node := &dictionaryTrieNode{
			id:       nextID,
			children: make(map[rune]*dictionaryTrieNode),
		}
		nextID++
		return node
	}

	root := newNode()
	dictionaryNames := make([]string, 0, len(dicts))
	for name := range dicts {
		dictionaryNames = append(dictionaryNames, name)
	}
	sort.Strings(dictionaryNames)

	for _, dictionaryName := range dictionaryNames {
		words := make([]string, 0, len(dicts[dictionaryName]))
		for word := range dicts[dictionaryName] {
			words = append(words, word)
		}
		sort.Strings(words)

		for _, word := range words {
			node := root
			for _, chr := range word {
				lower := unicode.ToLower(chr)
				if node.children[lower] == nil {
					node.children[lower] = newNode()
				}
				node = node.children[lower]
			}
			node.entries = append(node.entries, dictionaryTrieEntry{
				word:           word,
				rank:           dicts[dictionaryName][word],
				dictionaryName: dictionaryName,
			})
		}
	}

	return root
}

type l33tState struct {
	node *dictionaryTrieNode
	sub  map[rune]rune
}

func (state l33tState) clone() l33tState {
	return l33tState{
		node: state.node,
		sub:  copyRuneMap(state.sub),
	}
}

func (state l33tState) withSubstitution(l33t, plain rune) (l33tState, bool) {
	if existing, ok := state.sub[l33t]; ok && existing != plain {
		return l33tState{}, false
	}

	next := state.clone()
	if next.sub == nil {
		next.sub = make(map[rune]rune)
	}
	next.sub[l33t] = plain
	return next, true
}

func copyRuneMap(in map[rune]rune) map[rune]rune {
	if len(in) == 0 {
		return nil
	}
	out := make(map[rune]rune, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func reverseL33tTable(table map[string][]string) map[rune][]rune {
	reversed := make(map[rune][]rune)
	for plain, l33tChars := range table {
		plainRune := []rune(plain)[0]
		for _, l33tChr := range l33tChars {
			r := []rune(l33tChr)[0]
			reversed[r] = append(reversed[r], plainRune)
		}
	}

	for r := range reversed {
		slices.Sort(reversed[r])
	}
	return reversed
}

func stringSubstitution(sub map[rune]rune) map[string]string {
	res := make(map[string]string, len(sub))
	for l33t, plain := range sub {
		res[string(l33t)] = string(plain)
	}
	return res
}

func dedupL33tStates(states []l33tState) []l33tState {
	seen := make(map[string]bool)
	res := make([]l33tState, 0, len(states))

	for _, state := range states {
		key := l33tStateKey(state)
		if seen[key] {
			continue
		}
		seen[key] = true
		res = append(res, state)
	}
	return res
}

func l33tStateKey(state l33tState) string {
	var b strings.Builder
	b.WriteString(strconv.Itoa(state.node.id))
	b.WriteString("|")

	keys := make([]rune, 0, len(state.sub))
	for k := range state.sub {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for _, k := range keys {
		b.WriteRune(k)
		b.WriteRune('=')
		b.WriteRune(state.sub[k])
		b.WriteRune(';')
	}
	return b.String()
}
