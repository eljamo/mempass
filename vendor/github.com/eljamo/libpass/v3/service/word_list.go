package service

import (
	"errors"
	"fmt"

	"github.com/eljamo/libpass/v3/asset"
	"github.com/eljamo/libpass/v3/config"
)

type WordListService interface {
	GetWords() ([]string, error)
}

type DefaultWordListService struct {
	cfg      *config.Config
	rngSvc   RNGService
	wordList []string
}

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

func getWordList(wordList string, wordMinLength int, wordMaxLength int) ([]string, error) {
	if wordMaxLength < wordMinLength {
		return nil, fmt.Errorf("word_length_max (%d) must be greater than or equal to word_length_min (%d)", wordMaxLength, wordMinLength)
	}

	wl, err := asset.GetWordList(wordList)
	if err != nil {
		return nil, err
	}

	var fw []string
	for _, word := range wl {
		if len(word) >= wordMinLength && len(word) <= wordMaxLength {
			fw = append(fw, string(word))
		}
	}

	if len(fw) == 0 {
		return nil, fmt.Errorf("no words found in word list %s with a word_length_min of %d and word_length_max of %d", wordList, wordMinLength, wordMaxLength)
	}

	return fw, nil
}

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
