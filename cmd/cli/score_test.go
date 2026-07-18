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
		{0, "very weak"},
		{1, "weak"},
		{2, "fair"},
		{3, "strong"},
		{4, "very strong"},
		{-1, "unknown"},
		{5, "unknown"},
	}

	for _, tt := range tests {
		if got := scoreLabel(tt.score); got != tt.want {
			t.Errorf("scoreLabel(%d) = %q, want %q", tt.score, got, tt.want)
		}
	}
}

func TestFormatScoredPasswords(t *testing.T) {
	t.Parallel()

	pws := []string{"password", "!!12&paper&SEA&onto&12!!"}
	lines := formatScoredPasswords(pws)

	if len(lines) != 2 {
		t.Fatalf("formatScoredPasswords returned %d lines, want 2", len(lines))
	}

	// "password" is a top dictionary word, guaranteed weakest score
	if !strings.HasPrefix(lines[0], "password") || !strings.Contains(lines[0], "[0/4 very weak]") {
		t.Errorf("weak password line = %q, want prefix \"password\" and \"[0/4 very weak]\"", lines[0])
	}

	// generated-style password: 4 words, mixed case, digits, symbol padding
	if !strings.Contains(lines[1], "[4/4 very strong]") {
		t.Errorf("strong password line = %q, want \"[4/4 very strong]\"", lines[1])
	}

	// score columns align: '[' starts at the same index on every line
	if strings.Index(lines[0], "[") != strings.Index(lines[1], "[") {
		t.Errorf("score columns not aligned:\n%s\n%s", lines[0], lines[1])
	}
}

func TestFormatScoredPasswordsEmpty(t *testing.T) {
	t.Parallel()

	if lines := formatScoredPasswords(nil); len(lines) != 0 {
		t.Errorf("formatScoredPasswords(nil) returned %d lines, want 0", len(lines))
	}
}
