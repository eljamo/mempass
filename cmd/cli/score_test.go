package cli

import (
	"strings"
	"testing"
)

func TestScoreLabel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		score int
		want  string
	}{
		{0, "Very Weak"},
		{1, "Weak"},
		{2, "Fair"},
		{3, "Strong"},
		{4, "Very Strong"},
		{-1, "unknown"},
		{5, "unknown"},
	}

	for _, tt := range tests {
		if got := scoreLabel(tt.score); got != tt.want {
			t.Errorf("scoreLabel(%d) = %q, want %q", tt.score, got, tt.want)
		}
	}
}

func TestEvaluatePasswordsShowScore(t *testing.T) {
	t.Parallel()

	pws := []string{"password", "!!12&paper&SEA&onto&12!!"}
	lines, _ := evaluatePasswords(pws, true)

	if len(lines) != 2 {
		t.Fatalf("evaluatePasswords returned %d lines, want 2", len(lines))
	}

	// "password" is a top dictionary word, guaranteed weakest score
	if !strings.HasPrefix(lines[0], "password") || !strings.Contains(lines[0], "Unthrottled: [0/4, Very Weak]") {
		t.Errorf("weak password line = %q, want prefix \"password\" and \"Unthrottled: [0/4, Very Weak]\"", lines[0])
	}

	// generated-style password: 4 words, mixed case, digits, symbol padding
	if !strings.Contains(lines[1], "Throttled: [4/4, Very Strong]") || !strings.Contains(lines[1], "Unthrottled: [4/4, Very Strong]") {
		t.Errorf("strong password line = %q, want Throttled and Unthrottled both \"[4/4, Very Strong]\"", lines[1])
	}

	// score columns align: '(' starts at the same index on every line
	if strings.Index(lines[0], "(") != strings.Index(lines[1], "(") {
		t.Errorf("score columns not aligned:\n%s\n%s", lines[0], lines[1])
	}
}

func TestEvaluatePasswordsHideScore(t *testing.T) {
	t.Parallel()

	pws := []string{"password", "!!12&paper&SEA&onto&12!!"}
	lines, _ := evaluatePasswords(pws, false)

	if len(lines) != 2 || lines[0] != pws[0] || lines[1] != pws[1] {
		t.Errorf("evaluatePasswords(showScore=false) = %v, want unchanged %v", lines, pws)
	}
}

func TestEvaluatePasswordsWarning(t *testing.T) {
	t.Parallel()

	pws := []string{"password", "!!12&paper&SEA&onto&12!!"}
	_, warning := evaluatePasswords(pws, false)

	// the warning is produced regardless of showScore, so a weak
	// UnthrottledPasswordEntryScore is still surfaced without --score;
	// it's a single short line, not one per weak password
	if warning == "" {
		t.Fatal("evaluatePasswords warning = \"\", want a non-empty warning for the weak password")
	}
	if !strings.HasPrefix(warning, "WARNING: Crackable in") {
		t.Errorf("warning = %q, want it to start with \"WARNING: Crackable in\"", warning)
	}
}

func TestEvaluatePasswordsNoWarning(t *testing.T) {
	t.Parallel()

	pws := []string{"!!12&paper&SEA&onto&12!!"}
	_, warning := evaluatePasswords(pws, false)

	if warning != "" {
		t.Errorf("evaluatePasswords warning = %q, want \"\" (no weak passwords)", warning)
	}
}

func TestEvaluatePasswordsEmpty(t *testing.T) {
	t.Parallel()

	lines, warning := evaluatePasswords(nil, true)
	if len(lines) != 0 {
		t.Errorf("evaluatePasswords(nil) returned %d lines, want 0", len(lines))
	}
	if warning != "" {
		t.Errorf("evaluatePasswords(nil) warning = %q, want \"\"", warning)
	}
}
