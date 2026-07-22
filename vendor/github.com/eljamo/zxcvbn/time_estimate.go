package zxcvbn

import (
	"fmt"
	"math"
)

type EstimatedTimes struct {
	CrackTimesSeconds map[string]float64 `json:"crack_times_seconds"`
	CrackTimesDisplay map[string]string  `json:"crack_times_display"`
	Score             int                `json:"score"`

	// ThrottledPasswordEntryScore reflects resistance to an attacker whose
	// rate of password entry is limited (e.g. a login form that locks out
	// or slows down after failed attempts, such as online banking). It's
	// derived from online_throttling_100_per_hour.
	ThrottledPasswordEntryScore int `json:"throttled_password_entry_score"`

	// UnthrottledPasswordEntryScore reflects resistance to an attacker
	// whose rate of password entry is NOT limited - e.g. one who has
	// obtained the hash and is cracking it offline, with no lockout or
	// throttling to slow them down. It's the pessimistic score
	UnthrottledPasswordEntryScore int `json:"unthrottled_password_entry_score"`
}

// scoreTimeThresholdsSeconds are the crack-time boundaries (in seconds) between
// score bands 0-4. A given crack-time scenario is scored by how long it would
// actually take under that scenario, rather than by a scenario-agnostic guess
// count, so "score 3" means "about a year to crack" whether that's under a
// throttled online attack or an offline fast-hash crack.
var scoreTimeThresholdsSeconds = [4]float64{
	60 * 60,               // 1 hour
	60 * 60 * 24,          // 1 day
	60 * 60 * 24 * 31 * 3, // ~3 months
	60 * 60 * 24 * 365,    // 1 year
}

func scoreFromCrackTimeSeconds(seconds float64) int {
	for score, threshold := range scoreTimeThresholdsSeconds {
		if seconds < threshold {
			return score
		}
	}
	return len(scoreTimeThresholdsSeconds)
}

func guessesToScore(guesses float64) int {
	const DELTA = 5
	if guesses < 1e3+DELTA {
		// risky password: "too guessable"
		return 0
	}
	if guesses < 1e6+DELTA {
		// modest protection from throttled online attacks: "very guessable"
		return 1
	}
	if guesses < 1e8+DELTA {
		// modest protection from unthrottled online attacks: "somewhat guessable"
		return 2
	}
	if guesses < 1e10+DELTA {
		// modest protection from offline attacks: "safely unguessable"
		// assuming a salted, slow hash function like bcrypt, scrypt, PBKDF2, argon, etc
		return 3
	}
	// strong protection from offline attacks under same scenario: "very unguessable"
	return 4
}

func estimateAttackTimes(guesses float64) (t EstimatedTimes) {
	t.CrackTimesSeconds = make(map[string]float64)
	t.CrackTimesSeconds["online_throttling_100_per_hour"] = guesses * 3600.0 / 100.0
	t.CrackTimesSeconds["online_no_throttling_10_per_second"] = guesses / 10
	t.CrackTimesSeconds["offline_slow_hashing_1e4_per_second"] = guesses / 1e4
	t.CrackTimesSeconds["offline_fast_hashing_1e10_per_second"] = guesses / 1e10
	t.CrackTimesSeconds["offline_fast_hashing_1e13_per_second"] = guesses / 1e13

	t.CrackTimesDisplay = make(map[string]string)

	for scenario, seconds := range t.CrackTimesSeconds {
		t.CrackTimesDisplay[scenario] = displayTime(seconds)
	}

	t.Score = guessesToScore(guesses)
	t.ThrottledPasswordEntryScore = scoreFromCrackTimeSeconds(t.CrackTimesSeconds["online_throttling_100_per_hour"])
	t.UnthrottledPasswordEntryScore = scoreFromCrackTimeSeconds(t.CrackTimesSeconds["offline_fast_hashing_1e13_per_second"])
	return
}

func displayTime(seconds float64) string {
	minute := float64(60)
	hour := minute * 60
	day := hour * 24
	month := day * 31
	year := month * 12
	century := year * 100

	switch {
	case seconds < 1:
		return "less than a second"
	case seconds < minute:
		return strCount(seconds, "second")
	case seconds < hour:
		return strCount(seconds/minute, "minute")
	case seconds < day:
		return strCount(seconds/hour, "hour")
	case seconds < month:
		return strCount(seconds/day, "day")
	case seconds < year:
		return strCount(seconds/month, "month")
	case seconds < century:
		return strCount(seconds/year, "year")
	default:
		return "centuries"
	}
}

func strCount(count float64, base string) string {
	c := int(round(count, 0.5, 0))
	str := fmt.Sprintf("%d %s", c, base)
	if c > 1 {
		str += "s"
	}
	return str
}

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
