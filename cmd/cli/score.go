package cli

import (
	"fmt"
	"strings"

	"github.com/eljamo/zxcvbn"
)

// Labels for each zxcvbn score, indexed by score (0-4)
var scoreLabels = []string{"very weak", "weak", "fair", "strong", "very strong"}

const unknownScoreLabel = "unknown"

func scoreLabel(score int) string {
	if score < 0 || score >= len(scoreLabels) {
		return unknownScoreLabel
	}

	return scoreLabels[score]
}

func formatScoredPasswords(pws []string) []string {
	maxLen := 0
	for _, p := range pws {
		if len(p) > maxLen {
			maxLen = len(p)
		}
	}

	lines := make([]string, 0, len(pws))
	for _, p := range pws {
		r := zxcvbn.PasswordStrength(p, nil)
		pad := strings.Repeat(" ", maxLen-len(p))
		lines = append(lines, fmt.Sprintf("%s%s  [%d/4 %s]", p, pad, r.Score, scoreLabel(r.Score)))
	}

	return lines
}
