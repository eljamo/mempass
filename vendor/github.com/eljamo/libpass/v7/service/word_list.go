package service

import (
	"fmt"

	"github.com/eljamo/libpass/v7/asset"
	"github.com/eljamo/libpass/v7/config"
	"github.com/eljamo/libpass/v7/config/option"
)

// Defines the interface for a service that extracts words from word lists
type WordListService interface {
	// GetWords returns a slice of words extracted from a word list or an error
	// if the words cannot be extracted.
	GetWords() ([]string, error)
}

// Implements the interface WordListService, providing functionality to extract
// words from a word list.
type DefaultWordListService struct {
	cfg      *config.Settings
	rngSvc   RNGService
	wordList []string
}

const numWordMin = 2

// Creates a new instance of DefaultWordListService. It requires configuration
// and a random number generation service. It returns an error if the
// configuration is invalid.
func NewWordListService(cfg *config.Settings, rngSvc RNGService) (*DefaultWordListService, error) {
	if cfg.NumWords < numWordMin {
		return nil, fmt.Errorf("%s must be greater than or equal to %d", option.ConfigKeyNumWords, numWordMin)
	}

	wordList, err := getWordList(cfg.WordList, cfg.WordLengthMin, cfg.WordLengthMax)
	if err != nil {
		return nil, err
	}

	return &DefaultWordListService{
		cfg,
		rngSvc,
		wordList,
	}, nil
}

// Creates a word list based on provided criteria. It returns an error if the
// criteria are invalid or the word list cannot be created.
func getWordList(wordList string, wordMinLength int, wordMaxLength int) ([]string, error) {
	if wordMaxLength < wordMinLength {
		return nil, fmt.Errorf(
			"%s (%d) must be greater than or equal to %s (%d)",
			option.ConfigKeyWordLengthMax,
			wordMaxLength,
			option.ConfigKeyWordLengthMin,
			wordMinLength,
		)
	}

	wl, err := asset.GetFilteredWordList(wordList, wordMinLength, wordMaxLength)
	if err != nil {
		return nil, err
	}

	if len(wl) == 0 {
		return nil, fmt.Errorf(
			"no words found in %s (%s) with a %s of %d and %s of %d",
			option.ConfigKeyWordList,
			wordList,
			option.ConfigKeyWordLengthMin,
			wordMinLength,
			option.ConfigKeyWordLengthMax,
			wordMaxLength,
		)
	}

	return wl, nil
}

// Creates a slice of words randomly extracted from a word list. It returns
// an error if the slice cannot be created.
func (s *DefaultWordListService) GetWords() ([]string, error) {
	wll := len(s.wordList)
	wn, err := s.rngSvc.GenerateSliceWithMax(s.cfg.NumWords, wll)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random word slice index numbers: %w", err)
	}

	wl := make([]string, s.cfg.NumWords)
	for i, idx := range wn {
		if idx < 0 || idx >= wll {
			return nil, fmt.Errorf("number (%d) given out of range of word list (%d)", idx, wll)
		}
		wl[i] = s.wordList[idx]
	}

	return wl, nil
}
