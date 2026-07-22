package cli

import (
	"fmt"
	"math"
	"strings"

	"github.com/eljamo/zxcvbn"
)

// Labels for each zxcvbn score, indexed by score (0-4)
var scoreLabels = []string{"Very Weak", "Weak", "Fair", "Strong", "Very Strong"}

const unknownScoreLabel = "unknown"

// offlineFastHashingScenario is the pessimistic crack-time scenario key used
// to quote a concrete crack time in the weak-password warning below.
const offlineFastHashingScenario = "offline_fast_hashing_1e13_per_second"

// weakUnthrottledScoreThreshold is the UnthrottledPasswordEntryScore at or
// below which we warn, regardless of whether --score was passed. It's the
// pessimistic (offline, no-rate-limit) score, so a low value means the
// password would fall quickly to an attacker who has obtained the hash and
// isn't limited to a throttled login attempt rate - e.g. a leaked NTLM hash.
const weakUnthrottledScoreThreshold = 1

func scoreLabel(score int) string {
	if score < 0 || score >= len(scoreLabels) {
		return unknownScoreLabel
	}

	return scoreLabels[score]
}

// evaluatePasswords scores every generated password. When showScore is true,
// each line is annotated with the ThrottledPasswordEntryScore/
// UnthrottledPasswordEntryScore breakdown; otherwise lines are returned
// unchanged. Scoring itself always runs, regardless of showScore, so a
// single short warning is still produced for callers who never pass
// --score. A batch generated from one config tends to land in the same
// bucket, so rather than one line per weak password, only the single
// worst (fastest-to-crack) case across the whole batch is reported.
func evaluatePasswords(pws []string, showScore bool) (lines []string, warning string) {
	maxLen := 0
	for _, p := range pws {
		if len(p) > maxLen {
			maxLen = len(p)
		}
	}

	lines = make([]string, 0, len(pws))
	worstSeconds := math.Inf(1)
	var worstDisplay string

	for _, p := range pws {
		r := zxcvbn.PasswordStrength(p, nil)

		line := p
		if showScore {
			pad := strings.Repeat(" ", maxLen-len(p))
			line = fmt.Sprintf(
				"%s%s  (Throttled: [%d/4, %s], Unthrottled: [%d/4, %s])",
				p, pad,
				r.ThrottledPasswordEntryScore, scoreLabel(r.ThrottledPasswordEntryScore),
				r.UnthrottledPasswordEntryScore, scoreLabel(r.UnthrottledPasswordEntryScore),
			)
		}
		lines = append(lines, line)

		if r.UnthrottledPasswordEntryScore <= weakUnthrottledScoreThreshold {
			if secs := r.CrackTimesSeconds[offlineFastHashingScenario]; secs < worstSeconds {
				worstSeconds = secs
				worstDisplay = r.CrackTimesDisplay[offlineFastHashingScenario]
			}
		}
	}

	if worstDisplay != "" {
		fastest := worstDisplay
		if fastest != "less than a second" {
			fastest = "as little as " + fastest
		}
		warning = fmt.Sprintf("WARNING: Crackable in %s by an attacker with no rate limit", fastest)
	}

	return lines, warning
}
