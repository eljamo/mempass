package option

// Config key
const (
	ConfigKeyCaseTransform           string = "case_transform"
	ConfigKeyNumPasswords            string = "num_passwords"
	ConfigKeyNumWords                string = "num_words"
	ConfigKeyPaddingCharactersAfter  string = "padding_characters_after"
	ConfigKeyPaddingCharactersBefore string = "padding_characters_before"
	ConfigKeyPaddingCharacter        string = "padding_character"
	ConfigKeyPaddingDigitsAfter      string = "padding_digits_after"
	ConfigKeyPaddingDigitsBefore     string = "padding_digits_before"
	ConfigKeyPaddingType             string = "padding_type"
	ConfigKeyPadToLength             string = "pad_to_length"
	ConfigKeyPreset                  string = "preset"
	ConfigKeySeparatorAlphabet       string = "separator_alphabet"
	ConfigKeySeparatorCharacter      string = "separator_character"
	ConfigKeySymbolAlphabet          string = "symbol_alphabet"
	ConfigKeyWordLengthMax           string = "word_length_max"
	ConfigKeyWordLengthMin           string = "word_length_min"
	ConfigKeyWordList                string = "word_list"
)

// Word list constant
const (
	WordListAll           string = "ALL"
	WordListDoctorWho     string = "DOCTOR_WHO"
	WordListEN            string = "EN"
	WordListENSmall       string = "EN_SMALL"
	WordListGameOfThrones string = "GAME_OF_THRONES"
	WordListHarryPotter   string = "HARRY_POTTER"
	WordListMiddleEarth   string = "MIDDLE_EARTH"
	WordListPokemon       string = "POKEMON"
	WordListStarTrek      string = "STAR_TREK"
	WordListStarWars      string = "STAR_WARS"
)

// Preset constant
const (
	PresetAppleID       string = "APPLEID"
	PresetDefault       string = "DEFAULT"
	PresetNTLM          string = "NTLM"
	PresetSecurityQ     string = "SECURITYQ"
	PresetWeb16         string = "WEB16"
	PresetWeb16XKPasswd string = "WEB16_XKPASSWD"
	PresetWeb32         string = "WEB32"
	PresetWiFi          string = "WIFI"
	PresetXKCD          string = "XKCD"
	PresetXKCDXKPasswd  string = "XKCD_XKPASSWD"
)

// Case transform constant
const (
	CaseTransformAlternate                string = "ALTERNATE"
	CaseTransformAlternateLettercase      string = "ALTERNATE_LETTERCASE"
	CaseTransformCapitalise               string = "CAPITALISE"
	CaseTransformCapitaliseInvert         string = "CAPITALISE_INVERT"
	CaseTransformInvert                   string = "INVERT"
	CaseTransformLower                    string = "LOWER"
	CaseTransformLowerVowelUpperConsonant string = "LOWER_VOWEL_UPPER_CONSONANT"
	CaseTransformNone                     string = "NONE"
	CaseTransformRandom                   string = "RANDOM"
	CaseTransformSentence                 string = "SENTENCE"
	CaseTransformUpper                    string = "UPPER"
)

// Padding type constant
const (
	PaddingTypeAdaptive string = "ADAPTIVE"
	PaddingTypeFixed    string = "FIXED"
	PaddingTypeNone     string = "NONE"
)

const (
	PaddingCharacterRandom string = "RANDOM"
)

const (
	SeparatorCharacterRandom string = "RANDOM"
)
