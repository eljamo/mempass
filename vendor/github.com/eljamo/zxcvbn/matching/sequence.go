package matching

import (
	"regexp"
	"unicode/utf8"

	"github.com/eljamo/zxcvbn/match"
)

type sequenceMatch struct{}

const maxDelta = 5

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

var (
	reLower  = regexp.MustCompile(`^[a-z]+$`)
	reUpper  = regexp.MustCompile(`^[A-Z]+$`)
	reDigits = regexp.MustCompile(`^\d+$`)
)

func (sequenceMatch) Matches(password string) []*match.Match {
	matches := []*match.Match{}
	runes := []rune(password)
	if len(runes) <= 1 {
		return matches
	}

	// byteStart[k] is the byte offset of the k-th rune;
	// byteStart[len(runes)] == len(password).
	byteStart := make([]int, len(runes)+1)
	b := 0
	for k, r := range runes {
		byteStart[k] = b
		b += utf8.RuneLen(r)
	}
	byteStart[len(runes)] = b

	// i, j, delta are rune indices/deltas; Token, I, J are emitted in bytes.
	update := func(i, j, delta int) {
		absDelta := abs(delta)
		if j-i > 1 || absDelta == 1 {
			if absDelta > 0 && absDelta <= maxDelta {
				token := password[byteStart[i]:byteStart[j+1]]
				// conservatively stick with roman alphabet size.
				// (this could be improved)
				seqName := "unicode"
				seqSpace := 26
				switch {
				case reLower.MatchString(token):
					seqName = "lower"
				case reUpper.MatchString(token):
					seqName = "upper"
				case reDigits.MatchString(token):
					seqName = "digits"
					seqSpace = 10
				}
				matches = append(matches, &match.Match{
					Pattern:       "sequence",
					I:             byteStart[i],
					J:             byteStart[j+1] - 1,
					Token:         token,
					SequenceName:  seqName,
					SequenceSpace: seqSpace,
					Ascending:     delta > 0,
				})
			}
		}
	}

	i := 0
	lastDelta := 0 // null
	for k := 1; k <= len(runes)-1; k++ {
		delta := int(runes[k]) - int(runes[k-1])
		if k == 1 {
			lastDelta = delta
		}
		if delta == lastDelta {
			continue
		}
		j := k - 1
		update(i, j, lastDelta)
		i = j
		lastDelta = delta
	}

	update(i, len(runes)-1, lastDelta)
	return matches
}
