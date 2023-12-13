package option

// A slice of available presets
var Presets = []string{
	Default, AppleID, NTLM, SecurityQ, Web16, Web16XKPasswd, Web32, WiFi, XKCD,
	XKCDXKPasswd,
}

// A slice of special characters which can be used for padding and separator
// characters
var DefaultSpecialCharacters = []string{
	"!", "@", "$", "%", "^", "&", "*", "-", "+", "=", ":", "|", "~", "?", "/",
	".", ";",
}

// A slice of available options for padding
var PaddingTypes = []string{Adaptive, Fixed, None}

// A slice of available options for case transformation
var TransformTypes = []string{
	Alternate, AlternateLettercase, Capitalise, CapitaliseInvert, Invert, Lower,
	LowerVowelUpperConsonant, None, Random, Sentence, Upper,
}

// A slice of available options for padding and separator characters
var PaddingCharacterAndSeparatorCharacters = append([]string{Random}, DefaultSpecialCharacters...)

// A slice of available word lists
var WordLists = []string{
	All, DoctorWho, EN, ENSmall, GameOfThrones, HarryPotter, MiddleEarth,
	StarTrek, StarWars,
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
