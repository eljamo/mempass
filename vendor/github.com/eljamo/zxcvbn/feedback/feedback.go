package feedback

import (
	"regexp"
	"unicode/utf8"

	"github.com/eljamo/zxcvbn/match"
)

var (
	StartUpper = regexp.MustCompile(`^[A-Z][^A-Z]+$`)
	// AllUpper diverges from upstream on purpose; see GetFeedback's doc comment.
	AllUpper = regexp.MustCompile(`^[^a-z\d]+$`)
)

// Feedback represents the feedback for a weak password
type Feedback struct {
	Warning     string
	Suggestions []string
}

// GetFeedback follows
// https://github.com/dropbox/zxcvbn/blob/master/src/feedback.coffee
// with one deliberate divergence: AllUpper excludes digits (`^[^a-z\d]+$` vs
// upstream `^[^a-z]+$`), so digit-containing all-caps tokens like PASSWORD123
// skip the all-uppercase suggestion. Feedback phrasing is allowed to differ
// from scoring math; scoring keeps upstream semantics.
func GetFeedback(score int, sequence []*match.Match) Feedback {
	if len(sequence) == 0 {
		return Feedback{
			Warning: "",
			Suggestions: []string{
				"Use a few words, avoid common phrases",
				"No need for symbols, digits, or uppercase letters",
			},
		}
	}

	if score > 2 {
		return Feedback{}
	}

	longestMatch := sequence[0]
	for _, match := range sequence[1:] {
		if len(match.Token) > len(longestMatch.Token) {
			longestMatch = match
		}
	}

	feedback := getMatchFeedback(longestMatch, len(sequence) == 1)
	extraFeedback := "Add another word or two. Uncommon words are better"
	feedback.Suggestions = append(feedback.Suggestions, extraFeedback)

	return feedback
}

func getMatchFeedback(match *match.Match, isSoleMatch bool) Feedback {
	switch match.Pattern {
	case "dictionary":
		return getDictionaryMatchFeedback(match, isSoleMatch)
	case "spatial":
		warning := "Short keyboard patterns are easy to guess"
		if match.Turns == 1 {
			warning = "Straight rows of keys are easy to guess"
		}
		return Feedback{
			Warning:     warning,
			Suggestions: []string{"Use a longer keyboard pattern with more turns"},
		}
	case "repeat":
		warning := `Repeats like "abcabcabc" are only slightly harder to guess than "abc"`
		if utf8.RuneCountInString(match.BaseToken) == 1 {
			warning = `Repeats like "aaa" are easy to guess`
		}
		return Feedback{
			Warning:     warning,
			Suggestions: []string{"Avoid repeated words and characters."},
		}
	case "sequence":
		return Feedback{
			Warning:     "Sequences like \"abc\" or \"6543\" are easy to guess.",
			Suggestions: []string{"Avoid sequences"},
		}
	case "regex":
		if match.RegexName == "recent_year" {
			return Feedback{
				Warning: "Recent years are easy to guess.",
				Suggestions: []string{
					"Avoid recent years",
					"Avoid years that are associated with you",
				},
			}
		}
	case "date":
		return Feedback{
			Warning:     "Dates are often easy to guess.",
			Suggestions: []string{"Avoid dates and years that are associated with you."},
		}
	}

	return Feedback{}
}

func getDictionaryMatchFeedback(match *match.Match, isSoleMatch bool) Feedback {
	warning := ""
	switch match.DictionaryName {
	case "passwords":
		if isSoleMatch && !match.L33t && !match.Reversed {
			switch {
			case match.Rank <= 10:
				warning = "This is a top-10 common password"
			case match.Rank <= 100:
				warning = "This is a top-100 common password"
			default:
				warning = "This is a very common password"
			}
		} else if match.Guesses <= 10_000 {
			warning = "This is similar to a commonly used password"
		}
	case "english_wikipedia":
		if isSoleMatch {
			warning = "A word by itself is easy to guess"
		}
	case "surnames", "male_names", "female_names":
		if isSoleMatch {
			warning = "Names and surnames by themselves are easy to guess"
		} else {
			warning = "Common names and surnames are easy to guess"
		}
	}

	var suggestions []string
	word := match.Token

	if StartUpper.MatchString(word) {
		suggestions = append(suggestions, "Capitalization doesn't help very much")
	} else if AllUpper.MatchString(word) {
		suggestions = append(suggestions, "All-uppercase is almost as easy to guess as all-lowercase")
	}

	if match.Reversed && len(word) >= 4 {
		suggestions = append(suggestions, "Reversed words aren't much harder to guess")
	}
	if match.L33t {
		suggestions = append(suggestions, "Predictable substitutions like '@' instead of 'a' don't help very much")
	}

	return Feedback{
		Warning:     warning,
		Suggestions: suggestions,
	}
}
