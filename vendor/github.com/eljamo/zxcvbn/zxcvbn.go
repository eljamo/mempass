package zxcvbn

import (
	"math"
	"time"
	"unicode/utf8"

	"github.com/eljamo/zxcvbn/feedback"
	"github.com/eljamo/zxcvbn/match"
	"github.com/eljamo/zxcvbn/matching"
	"github.com/eljamo/zxcvbn/scoring"
)

type Result struct {
	Guesses      float64           `json:"guesses"`
	GuessesLog10 float64           `json:"guesses_log10"`
	Sequence     []*match.Match    `json:"sequence"`
	CalcTime     float64           `json:"calc_time"`
	Feedback     feedback.Feedback `json:"feedback"`
	EstimatedTimes
}

type Config struct {
	CustomDictionaries map[string][]string
}

type Estimator struct {
	matcher matching.Omnimatcher
}

func NewEstimator(config Config) *Estimator {
	return &Estimator{
		matcher: matching.NewOmnimatcher(config.CustomDictionaries),
	}
}

var defaultEstimator = NewEstimator(Config{})

func PasswordStrength(password string, userInputs []string) Result {
	return defaultEstimator.PasswordStrength(password, userInputs)
}

func (e *Estimator) PasswordStrength(password string, userInputs []string) Result {
	start := time.Now()
	var result Result
	if !utf8.ValidString(password) {
		// Do not evaluate passwords containing invalid utf8
		// => those will be reported as weak passwords
		return result
	}
	matches := e.matcher.Omnimatch(password, userInputs)
	seq := scoring.MostGuessableMatchSequence(password, matches, false)
	result.Sequence = seq.Sequence
	result.Guesses = seq.Guesses
	result.GuessesLog10 = math.Log10(seq.Guesses)
	result.EstimatedTimes = estimateAttackTimes(seq.Guesses)
	result.Feedback = feedback.GetFeedback(result.Score, result.Sequence)
	result.CalcTime = time.Since(start).Seconds()
	return result
}
