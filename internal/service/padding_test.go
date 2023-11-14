package service

import (
	"reflect"
	"strings"
	"testing"

	"github.com/eljamo/mempass/internal/config"
)

func TestPad(t *testing.T) {
	cfg := &config.Config{
		PaddingDigitsBefore:     2,
		PaddingDigitsAfter:      2,
		SeparatorCharacter:      "-",
		PaddingType:             config.FIXED,
		PaddingCharactersBefore: 2,
		PaddingCharactersAfter:  2,
		PadToLength:             20,
		PaddingCharacter:        "*",
		SymbolAlphabet:          []string{"!"},
	}
	rngs := NewMockRNGService()
	paddingService := NewPaddingService(cfg, rngs)

	input := []string{"-test-"}
	expected := "**11-test-11**"

	result, err := paddingService.Pad(input)
	if err != nil {
		t.Errorf("Pad() error = %v, expectErr %v", err, false)
	}
	if result != expected {
		t.Errorf("Pad() = %v, expected %v", result, expected)
	}
}

func TestDigits(t *testing.T) {
	rngs := NewMockRNGService()
	tests := []struct {
		name      string
		input     []string
		before    int
		after     int
		expected  []string
		expectErr bool
	}{
		{
			name:      "valid padding",
			input:     []string{"a", "b", "c"},
			expected:  []string{"1", "1", "a", "b", "c", "1", "1"},
			before:    2,
			after:     2,
			expectErr: false,
		},
		{
			name:      "no padding",
			input:     []string{"a", "b", "c"},
			expected:  []string{"a", "b", "c"},
			before:    0,
			after:     0,
			expectErr: false,
		},
		{
			name:      "negative padding before",
			input:     []string{"a", "b", "c"},
			before:    -1,
			after:     0,
			expectErr: true,
		},
		{
			name:      "negative padding after",
			input:     []string{"a", "b", "c"},
			before:    0,
			after:     -1,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				PaddingDigitsBefore: tt.before,
				PaddingDigitsAfter:  tt.after,
			}
			s := NewPaddingService(cfg, rngs)
			got, err := s.digits(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("digits() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("digits() got = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestGenerateRandomDigits(t *testing.T) {
	cfg := &config.Config{}
	tests := []struct {
		name     string
		count    int
		number   int
		expected []string
	}{
		{
			name:     "generate 3 digits",
			count:    3,
			number:   1,
			expected: []string{"1", "1", "1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rngs := NewMockRNGService()
			s := NewPaddingService(cfg, rngs)

			got, err := s.generateRandomDigits(tt.count)
			if tt.count < 0 {
				if err == nil {
					t.Errorf("generateRandomDigits() expected error for count %d, got nil", tt.count)
				}
			} else {
				if err != nil {
					t.Errorf("generateRandomDigits() error = %v, expectErr nil", err)
				}
				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("generateRandomDigits() got = %v, expected %v", got, tt.expected)
				}
			}
		})
	}
}

func TestRemoveEdgeSeparatorCharacter(t *testing.T) {
	cfg := &config.Config{SeparatorCharacter: "-"}

	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no separator at edges",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "separator at start",
			input:    []string{"-", "a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "separator at end",
			input:    []string{"a", "b", "c", "-"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "separator at both ends",
			input:    []string{"-", "a", "b", "c", "-"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "input with only separators",
			input:    []string{"-", "-"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DefaultPaddingService{
				cfg: cfg,
			}
			got := s.removeEdgeSeparatorCharacter(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("removeEdgeSeparatorCharacter() got = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestSymbols(t *testing.T) {
	rngs := NewMockRNGService()
	tests := []struct {
		name       string
		cfg        *config.Config
		pw         string
		want       string
		expectErr  bool
		errMessage string
	}{
		{
			name: "fixed padding with specific character",
			cfg: &config.Config{
				PaddingType:             config.FIXED,
				PaddingCharacter:        "*",
				PaddingCharactersBefore: 2,
				PaddingCharactersAfter:  2,
				PadToLength:             10,
			},
			pw:   "pass",
			want: "**pass**",
		},
		{
			name: "adaptive padding to specific length",
			cfg: &config.Config{
				PaddingType:      config.ADAPTIVE,
				PaddingCharacter: "*",
				PadToLength:      10,
			},
			pw:   "pass",
			want: "pass******",
		},
		{
			name: "no padding",
			cfg: &config.Config{
				PaddingType: config.NONE,
			},
			pw:   "password",
			want: "password",
		},
		{
			name: "random padding character",
			cfg: &config.Config{
				PaddingType:             config.FIXED,
				PaddingCharacter:        config.RANDOM,
				PaddingCharactersBefore: 2,
				PaddingCharactersAfter:  2,
				PadToLength:             10,
				SymbolAlphabet:          []string{"!", "-", "_"},
			},
			pw:   "pass",
			want: "--pass--",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPaddingService(tt.cfg, rngs)
			got, err := s.symbols(tt.pw)
			if (err != nil) != tt.expectErr {
				t.Errorf("symbols() error = %v, expectErr %v", err, tt.expectErr)
			}
			if tt.expectErr && err.Error() != tt.errMessage {
				t.Errorf("symbols() error = %v, expectErr message %v", err, tt.errMessage)
			}
			if !tt.expectErr && got != tt.want {
				t.Errorf("symbols() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFixed(t *testing.T) {
	rngs := NewMockRNGService()
	tests := []struct {
		name        string
		pw          string
		char        string
		before      int
		after       int
		want        string
		expectErr   bool
		errContains string
	}{
		{
			name:   "equal padding on both sides",
			pw:     "abc",
			char:   "*",
			before: 2,
			after:  2,
			want:   "**abc**",
		},
		{
			name:   "no padding",
			pw:     "abc",
			char:   "*",
			before: 0,
			after:  0,
			want:   "abc",
		},
		{
			name:   "padding before only",
			pw:     "abc",
			char:   "#",
			before: 3,
			after:  0,
			want:   "###abc",
		},
		{
			name:   "padding after only",
			pw:     "abc",
			char:   "+",
			before: 0,
			after:  4,
			want:   "abc++++",
		},
		{
			name:        "negative padding before",
			pw:          "abc",
			char:        "*",
			before:      -1,
			after:       2,
			expectErr:   true,
			errContains: "padding_characters_before and padding_characters_after must be greater than or equal to 0",
		},
		{
			name:        "negative padding after",
			pw:          "abc",
			char:        "*",
			before:      1,
			after:       -2,
			expectErr:   true,
			errContains: "padding_characters_before and padding_characters_after must be greater than or equal to 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				PaddingDigitsBefore: tt.before,
				PaddingDigitsAfter:  tt.after,
			}
			s := NewPaddingService(cfg, rngs)
			got, err := s.fixed(tt.pw, tt.char, tt.before, tt.after)
			if (err != nil) != tt.expectErr {
				t.Fatalf("fixed() error = %v, expectErr %v", err, tt.expectErr)
			}
			if tt.expectErr {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("fixed() error = %v, want to contain %s", err, tt.errContains)
				}
			} else {
				if got != tt.want {
					t.Errorf("fixed() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestAdaptive(t *testing.T) {
	cfg := &config.Config{
		PadToLength: 10,
	}

	tests := []struct {
		name       string
		pw         string
		char       string
		padLen     int
		want       string
		expectErr  bool
		errMessage string
	}{
		{
			name:      "no padding needed",
			pw:        "1234567890",
			char:      "*",
			padLen:    10,
			want:      "1234567890",
			expectErr: false,
		},
		{
			name:      "padding needed",
			pw:        "12345",
			char:      "*",
			padLen:    10,
			want:      "12345*****",
			expectErr: false,
		},
		{
			name:      "empty password",
			pw:        "",
			char:      "*",
			padLen:    10,
			want:      "**********",
			expectErr: false,
		},
		{
			name:       "negative pad length",
			pw:         "12345",
			char:       "*",
			padLen:     -1,
			expectErr:  true,
			errMessage: "pad_to_length must be greater than or equal to 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg.PadToLength = tt.padLen
			s := &DefaultPaddingService{
				cfg: cfg,
			}
			got, err := s.adaptive(tt.pw, tt.char, tt.padLen)
			if (err != nil) != tt.expectErr {
				t.Errorf("adaptive() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if tt.expectErr && err.Error() != tt.errMessage {
				t.Errorf("adaptive() error = %v, expectErr message %v", err, tt.errMessage)
			}
			if !tt.expectErr && got != tt.want {
				t.Errorf("adaptive() got = %v, want %v", got, tt.want)
			}
		})
	}
}
