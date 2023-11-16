package config

const (
	ALL             string = "ALL"
	DOCTOR_WHO      string = "DOCTOR_WHO"
	EN              string = "EN"
	EN_SMALL        string = "EN_SMALL"
	GAME_OF_THRONES string = "GAME_OF_THRONES"
	HARRY_POTTER    string = "HARRY_POTTER"
	MIDDLE_EARTH    string = "MIDDLE_EARTH"
	POKEMON         string = "POKEMON"
	STAR_TREK       string = "STAR_TREK"
	STAR_WARS       string = "STAR_WARS"
)

const (
	APPLEID        string = "APPLEID"
	DEFAULT        string = "DEFAULT"
	NTLM           string = "NTLM"
	SECURITYQ      string = "SECURITYQ"
	WEB16          string = "WEB16"
	WEB16_XKPASSWD string = "WEB16_XKPASSWD"
	WEB32          string = "WEB32"
	WIFI           string = "WIFI"
	XKCD           string = "XKCD"
	XKCD_XKPASSWD  string = "XKCD_XKPASSWD"
)

const (
	ADAPTIVE   string = "ADAPTIVE"
	ALTERNATE  string = "ALTERNATE"
	CAPITALISE string = "CAPITALISE"
	FIXED      string = "FIXED"
	INVERT     string = "INVERT"
	LOWER      string = "LOWER"
	NONE       string = "NONE"
	RANDOM     string = "RANDOM"
	SPECIFIED  string = "SPECIFIED"
	UPPER      string = "UPPER"
)

var (
	Preset = []string{
		DEFAULT, APPLEID, NTLM, SECURITYQ, WEB16, WEB16_XKPASSWD, WEB32, WIFI,
		XKCD, XKCD_XKPASSWD,
	}
	DefaultSpecialCharacters = []string{
		"!", "@", "$", "%", "^", "&", "*", "-", "-", "+", "=", ":", "|", "~",
		"?", "/", ".", ";",
	}
	PaddingType = []string{
		ADAPTIVE, FIXED, NONE,
	}
	SeparatorCharacterType = []string{
		NONE, RANDOM, SPECIFIED,
	}
	TransformType = []string{
		ALTERNATE, CAPITALISE, INVERT, LOWER, NONE, RANDOM, UPPER,
	}
	PaddingCharacterAndSeparatorCharacter = append(
		[]string{RANDOM}, DefaultSpecialCharacters...,
	)
	WordLists = []string{
		ALL, DOCTOR_WHO, EN, EN_SMALL, GAME_OF_THRONES, HARRY_POTTER, MIDDLE_EARTH,
		STAR_TREK, STAR_WARS,
	}
)

const (
	PRESET_KEY    string = "preset"
	WORD_LIST_KEY string = "word_list"
)

type Config struct {
	CaseTransform           string   `json:"case_transform,omitempty"`
	NumPasswords            int      `json:"num_passwords,omitempty"`
	NumWords                int      `json:"num_words,omitempty"`
	PaddingCharactersAfter  int      `json:"padding_characters_after,omitempty"`
	PaddingCharactersBefore int      `json:"padding_characters_before,omitempty"`
	PaddingCharacter        string   `json:"padding_character,omitempty"`
	PaddingDigitsAfter      int      `json:"padding_digits_after,omitempty"`
	PaddingDigitsBefore     int      `json:"padding_digits_before,omitempty"`
	PaddingType             string   `json:"padding_type,omitempty"`
	PadToLength             int      `json:"pad_to_length,omitempty"`
	Preset                  string   `json:"preset,omitempty"`
	SeparatorAlphabet       []string `json:"separator_alphabet,omitempty"`
	SeparatorCharacter      string   `json:"separator_character,omitempty"`
	SymbolAlphabet          []string `json:"symbol_alphabet,omitempty"`
	WordLengthMax           int      `json:"word_length_max,omitempty"`
	WordLengthMin           int      `json:"word_length_min,omitempty"`
	WordList                string   `json:"word_list,omitempty"`
}

var WordListDescriptionMap = map[string]string{
	ALL:             "A combination of all word lists (44900+ words)",
	DOCTOR_WHO:      "A Doctor Who word list (11300+ words)",
	EN:              "A list of English words (14900+ words)",
	EN_SMALL:        "A small list of English words (8600+ words)",
	GAME_OF_THRONES: "A Game of Thrones word list (8200+ words)",
	HARRY_POTTER:    "A Harry Potter word list (12500+ words)",
	MIDDLE_EARTH:    "A Middle Earth word list containing words from The Hobbit, Lord of the Rings, The Silmarillion, and more (15400+ words)",
	POKEMON:         "A Pokemon word list (9000+ words)",
	STAR_TREK:       "A Star Trek word list (8000+ words)",
	STAR_WARS:       "A Star Wars word list (12000+ words)",
}

var PresetDescriptionMap = map[string]string{
	APPLEID:        "A preset respecting the many prerequisites Apple places on Apple ID passwords. The preset also limits itself to symbols found on the iOS letter and number keyboards (i.e. not the awkward to reach symbol keyboard)",
	DEFAULT:        "The default preset resulting in a password consisting of 3 random words of between 4 and 8 letters with alternating case separated by a random character, with two random digits before and after, and padded with two random characters front and back",
	NTLM:           "A preset for 14 character Windows NTLMv1 password. WARNING - only use this preset if you have to, it is too short to be acceptably secure",
	SECURITYQ:      "A preset for creating fake answers to security questions",
	WEB16:          "A preset for websites that insist passwords not be longer than 16 characters",
	WEB16_XKPASSWD: "A preset for websites that insist passwords not be longer than 16 characters, the same as the one found on xkpasswd.net.",
	WEB32:          "A preset for websites that allow passwords up to 32 characteres long",
	WIFI:           "A preset for generating 63 character long WPA2 keys",
	XKCD:           "A preset for generating passwords similar to the example in the original XKCD cartoon, but with a dash to separate the four randomly capitalised words, two digits and a random special characters.",
	XKCD_XKPASSWD:  "A preset for generating passwords similar to the example in the original XKCD cartoon, but with a dash to separate the four random words, and the capitalisation randomised to add sufficient entropy to avoid warnings.",
}
