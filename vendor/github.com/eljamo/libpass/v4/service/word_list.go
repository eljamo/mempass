package service

import (
	"errors"
	"fmt"

	"github.com/eljamo/libpass/v4/asset"
	"github.com/eljamo/libpass/v4/config"
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
	cfg      *config.Config
	rngSvc   RNGService
	wordList []string
}

// Creates a new instance of DefaultWordListService. It requires configuration
// and a random number generation service. It returns an error if the
// configuration is invalid.
func NewWordListService(cfg *config.Config, rngSvc RNGService) (*DefaultWordListService, error) {
	if cfg.NumWords < 2 {
		return nil, errors.New("num_words must be greater than or equal to 2")
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
		return nil, fmt.Errorf("word_length_max (%d) must be greater than or equal to word_length_min (%d)", wordMaxLength, wordMinLength)
	}

	wl, err := asset.GetFilteredWordList(wordList, wordMinLength, wordMaxLength)
	if err != nil {
		return nil, err
	}

	if len(wl) == 0 {
		return nil, fmt.Errorf("no words found in word list %s with a word_length_min of %d and word_length_max of %d", wordList, wordMinLength, wordMaxLength)
	}

	return wl, nil
}

// Creates a slice of words randomly extracted from a word list. It returns
// an error if the slice cannot be created.
func (s *DefaultWordListService) GetWords() ([]string, error) {
	wll := len(s.wordList)
	wn, err := s.rngSvc.GenerateSliceWithMax(s.cfg.NumWords, wll)
	if err != nil {
		return nil, err
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
