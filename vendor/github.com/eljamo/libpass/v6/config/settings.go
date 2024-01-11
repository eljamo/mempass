package config

import "github.com/eljamo/libpass/v6/config/option"

type Settings struct {
	// The type of case transformation to apply to the words
	CaseTransform string `key:"case_transform" json:"case_transform,omitempty"`
	// The number of passwords to generate
	NumPasswords int `key:"num_passwords" json:"num_passwords,omitempty"`
	// The number of words to use in the password
	NumWords int `key:"num_words" json:"num_words,omitempty"`
	// The number of padding characters to add after the password
	PaddingCharactersAfter int `key:"padding_characters_after" json:"padding_characters_after,omitempty"`
	// The number of padding characters to add before the password
	PaddingCharactersBefore int `key:"padding_characters_before" json:"padding_characters_before,omitempty"`
	// The character to use for padding
	PaddingCharacter string `key:"padding_character" json:"padding_character,omitempty"`
	// Te number of padding digits to add after the password
	PaddingDigitsAfter int `key:"padding_digits_after" json:"padding_digits_after,omitempty"`
	// The number of padding digits to add before the password
	PaddingDigitsBefore int `key:"padding_digits_before" json:"padding_digits_before,omitempty"`
	// The type of padding to apply to the password
	PaddingType string `key:"padding_type" json:"padding_type,omitempty"`
	// The length to pad the password to
	PadToLength int `key:"pad_to_length" json:"pad_to_length,omitempty"`
	// The preset to use for generating the password
	Preset string `key:"preset" json:"preset,omitempty"`
	// The alphabet to use for the separator character when using a random character
	SeparatorAlphabet []string `key:"separator_alphabet" json:"separator_alphabet,omitempty"`
	// The character to use to separate the words
	SeparatorCharacter string `key:"separator_character" json:"separator_character,omitempty"`
	// The alphabet to use for the symbol padding character when random
	SymbolAlphabet []string `key:"symbol_alphabet" json:"symbol_alphabet,omitempty"`
	// The maximum length of a word to use in the password
	WordLengthMax int `key:"word_length_max" json:"word_length_max,omitempty"`
	// The minimum length of a word to use in the password
	WordLengthMin int `key:"word_length_min" json:"word_length_min,omitempty"`
	// The word list to use for generating the password
	WordList string `key:"word_list" json:"word_list,omitempty"`
}

// DefaultSettings returns a new Settings struct with the default values set.
// This is used when no settings are given to the New function.
func DefaultSettings() *Settings {
	return &Settings{
		CaseTransform:           option.Random,
		NumPasswords:            3,
		NumWords:                3,
		PaddingCharacter:        option.Random,
		PaddingCharactersAfter:  2,
		PaddingCharactersBefore: 2,
		PaddingDigitsAfter:      2,
		PaddingDigitsBefore:     2,
		PaddingType:             option.Fixed,
		SeparatorAlphabet:       option.DefaultSpecialCharacters,
		SeparatorCharacter:      option.Random,
		SymbolAlphabet:          option.DefaultSpecialCharacters,
		WordLengthMax:           8,
		WordLengthMin:           4,
		WordList:                option.EN,
	}
}
