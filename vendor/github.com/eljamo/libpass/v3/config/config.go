package config

const (
	PresetKey   = "preset"
	WordListKey = "word_list"
)

const (
	All           = "ALL"
	DoctorWho     = "DOCTOR_WHO"
	EN            = "EN"
	ENSmall       = "EN_SMALL"
	GameOfThrones = "GAME_OF_THRONES"
	HarryPotter   = "HARRY_POTTER"
	MiddleEarth   = "MIDDLE_EARTH"
	Pokemon       = "POKEMON"
	StarTrek      = "STAR_TREK"
	StarWars      = "STAR_WARS"
)

const (
	AppleID       = "APPLEID"
	Default       = "DEFAULT"
	NTLM          = "NTLM"
	SecurityQ     = "SECURITYQ"
	Web16         = "WEB16"
	Web16XKPasswd = "WEB16_XKPASSWD"
	Web32         = "WEB32"
	WiFi          = "WIFI"
	XKCD          = "XKCD"
	XKCDXKPasswd  = "XKCD_XKPASSWD"
)

const (
	None   = "NONE"
	Random = "RANDOM"
)

const (
	Alternate                = "ALTERNATE"
	AlternateLettercase      = "ALTERNATE_LETTERCASE"
	Capitalise               = "CAPITALISE"
	CapitaliseInvert         = "CAPITALISE_INVERT"
	Invert                   = "INVERT" // Same as CapitaliseInvert but reserved to maintain compatibility with xkpasswd.net generated configs
	Lower                    = "LOWER"
	LowerVowelUpperConsonant = "LOWER_VOWEL_UPPER_CONSONANT"
	Sentence                 = "SENTENCE"
	Upper                    = "UPPER"
)

const (
	Adaptive = "ADAPTIVE"
	Fixed    = "FIXED"
)

var Preset = []string{
	Default, AppleID, NTLM, SecurityQ, Web16, Web16XKPasswd, Web32, WiFi, XKCD,
	XKCDXKPasswd,
}

var DefaultSpecialCharacters = []string{
	"!", "@", "$", "%", "^", "&", "*", "-", "+", "=", ":", "|", "~", "?", "/", ".", ";",
}

var PaddingType = []string{Adaptive, Fixed, None}

var TransformType = []string{
	Alternate, AlternateLettercase, Capitalise, CapitaliseInvert, Invert, Lower,
	LowerVowelUpperConsonant, None, Random, Sentence, Upper,
}

var PaddingCharacterAndSeparatorCharacter = append([]string{Random}, DefaultSpecialCharacters...)

var WordLists = []string{
	All, DoctorWho, EN, ENSmall, GameOfThrones, HarryPotter, MiddleEarth,
	StarTrek, StarWars,
}

type Config struct {
	CaseTransform           string   `key:"case_transform" json:"case_transform,omitempty"`
	NumPasswords            int      `key:"num_passwords" json:"num_passwords,omitempty"`
	NumWords                int      `key:"num_words" json:"num_words,omitempty"`
	PaddingCharactersAfter  int      `key:"padding_characters_after" json:"padding_characters_after,omitempty"`
	PaddingCharactersBefore int      `key:"padding_characters_before" json:"padding_characters_before,omitempty"`
	PaddingCharacter        string   `key:"padding_character" json:"padding_character,omitempty"`
	PaddingDigitsAfter      int      `key:"padding_digits_after" json:"padding_digits_after,omitempty"`
	PaddingDigitsBefore     int      `key:"padding_digits_before" json:"padding_digits_before,omitempty"`
	PaddingType             string   `key:"padding_type" json:"padding_type,omitempty"`
	PadToLength             int      `key:"pad_to_length" json:"pad_to_length,omitempty"`
	Preset                  string   `key:"preset" json:"preset,omitempty"`
	SeparatorAlphabet       []string `key:"separator_alphabet" json:"separator_alphabet,omitempty"`
	SeparatorCharacter      string   `key:"separator_character" json:"separator_character,omitempty"`
	SymbolAlphabet          []string `key:"symbol_alphabet" json:"symbol_alphabet,omitempty"`
	WordLengthMax           int      `key:"word_length_max" json:"word_length_max,omitempty"`
	WordLengthMin           int      `key:"word_length_min" json:"word_length_min,omitempty"`
	WordList                string   `key:"word_list" json:"word_list,omitempty"`
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
