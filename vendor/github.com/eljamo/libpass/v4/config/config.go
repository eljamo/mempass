package config

// Config key
const (
	PresetKey   string = "preset"
	WordListKey string = "word_list"
)

// Word list constant
const (
	All           string = "ALL"
	DoctorWho     string = "DOCTOR_WHO"
	EN            string = "EN"
	ENSmall       string = "EN_SMALL"
	GameOfThrones string = "GAME_OF_THRONES"
	HarryPotter   string = "HARRY_POTTER"
	MiddleEarth   string = "MIDDLE_EARTH"
	Pokemon       string = "POKEMON"
	StarTrek      string = "STAR_TREK"
	StarWars      string = "STAR_WARS"
)

// Preset constant
const (
	AppleID       string = "APPLEID"
	Default       string = "DEFAULT"
	NTLM          string = "NTLM"
	SecurityQ     string = "SECURITYQ"
	Web16         string = "WEB16"
	Web16XKPasswd string = "WEB16_XKPASSWD"
	Web32         string = "WEB32"
	WiFi          string = "WIFI"
	XKCD          string = "XKCD"
	XKCDXKPasswd  string = "XKCD_XKPASSWD"
)

// Shared constant
const (
	None   string = "NONE"
	Random string = "RANDOM"
)

// Case transform constant
const (
	Alternate           string = "ALTERNATE"
	AlternateLettercase string = "ALTERNATE_LETTERCASE"
	Capitalise          string = "CAPITALISE"
	CapitaliseInvert    string = "CAPITALISE_INVERT"
	// The same as CapitaliseInvert but reserved to maintain compatibility with xkpasswd.net generated configs
	Invert                   string = "INVERT"
	Lower                    string = "LOWER"
	LowerVowelUpperConsonant string = "LOWER_VOWEL_UPPER_CONSONANT"
	Sentence                 string = "SENTENCE"
	Upper                    string = "UPPER"
)

// Padding type constant
const (
	Adaptive string = "ADAPTIVE"
	Fixed    string = "FIXED"
)

// A slice of available presets
var Preset = []string{
	Default, AppleID, NTLM, SecurityQ, Web16, Web16XKPasswd, Web32, WiFi, XKCD,
	XKCDXKPasswd,
}

// A slice of special characters which can be used for padding and separator
// characters
var DefaultSpecialCharacters = []string{
	"!", "@", "$", "%", "^", "&", "*", "-", "+", "=", ":", "|", "~", "?", "/", ".", ";",
}

// A slice of available options for padding
var PaddingType = []string{Adaptive, Fixed, None}

// A slice of available options for case transformation
var TransformType = []string{
	Alternate, AlternateLettercase, Capitalise, CapitaliseInvert, Invert, Lower,
	LowerVowelUpperConsonant, None, Random, Sentence, Upper,
}

// A slice of available options for padding and separator characters
var PaddingCharacterAndSeparatorCharacter = append([]string{Random}, DefaultSpecialCharacters...)

// A slice of available word lists
var WordLists = []string{
	All, DoctorWho, EN, ENSmall, GameOfThrones, HarryPotter, MiddleEarth,
	StarTrek, StarWars,
}

type Config struct {
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

var WordListDescriptionMap = map[string]string{
	All:           "A combination of all word lists (44900+ words)",
	DoctorWho:     "A Doctor Who word list (11300+ words)",
	EN:            "A list of English words (14900+ words)",
	ENSmall:       "A small list of English words (8600+ words)",
	GameOfThrones: "A Game of Thrones word list (8200+ words)",
	HarryPotter:   "A Harry Potter word list (12500+ words)",
	MiddleEarth:   "A Middle Earth word list containing words from The Hobbit, Lord of the Rings, The Silmarillion, and more (15400+ words)",
	Pokemon:       "A Pokemon word list (9000+ words)",
	StarTrek:      "A Star Trek word list (8000+ words)",
	StarWars:      "A Star Wars word list (12000+ words)",
}

var PresetDescriptionMap = map[string]string{
	AppleID:       "A preset respecting the many prerequisites Apple places on Apple ID passwords. The preset also limits itself to symbols found on the iOS letter and number keyboards (i.e. not the awkward to reach symbol keyboard)",
	Default:       "The default preset resulting in a password consisting of 3 random words of between 4 and 8 letters with alternating case separated by a random character, with two random digits before and after, and padded with two random characters front and back",
	NTLM:          "A preset for 14 character Windows NTLMv1 password. WARNING - only use this preset if you have to, it is too short to be acceptably secure",
	SecurityQ:     "A preset for creating fake answers to security questions",
	Web16:         "A preset for websites that insist passwords not be longer than 16 characters",
	Web16XKPasswd: "A preset for websites that insist passwords not be longer than 16 characters, the same as the one found on xkpasswd.net.",
	Web32:         "A preset for websites that allow passwords up to 32 characteres long",
	WiFi:          "A preset for generating 63 character long WPA2 keys",
	XKCD:          "A preset for generating passwords similar to the example in the original XKCD cartoon, but with a dash to separate the four randomly capitalised words, two digits and a random special characters.",
	XKCDXKPasswd:  "A preset for generating passwords similar to the example in the original XKCD cartoon, but with a dash to separate the four random words, and the capitalisation randomised to add sufficient entropy to avoid warnings.",
}
