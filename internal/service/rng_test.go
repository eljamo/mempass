package service

import (
	"testing"
)

func TestGenerateN(t *testing.T) {
	tests := []struct {
		max       int
		expectErr bool
	}{
		{max: 100, expectErr: false},
		{max: -1, expectErr: true},
	}

	rng := NewRNGService()

	for _, tc := range tests {
		generated, err := rng.GenerateN(tc.max)

		if tc.expectErr {
			if err == nil {
				t.Errorf("Expected an error for max = %v, but got none", tc.max)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for max = %v: %v", tc.max, err)
			}

			if generated < 0 || generated >= tc.max {
				t.Errorf("Generated number is out of bounds for max = %v: got %v", tc.max, generated)
			}
		}
	}
}

func TestGenerate(t *testing.T) {
	rng := NewRNGService()
	generated, err := rng.Generate()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if generated < 0 {
		t.Errorf("Generated number is negative: got %v", generated)
	}
}
