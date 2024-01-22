package option

// A slice of available presets
var Presets = []string{
	PresetDefault, PresetAppleID, PresetNTLM, PresetSecurityQ, PresetWeb16,
	PresetWeb16XKPasswd, PresetWeb32, PresetWiFi, PresetXKCD, PresetXKCDXKPasswd,
}

// A slice of special characters which can be used for padding and separator
// characters
var DefaultSpecialCharacters = []string{
	"!", "@", "$", "%", "^", "&", "*", "-", "+", "=", ":", "|", "~", "?", "/",
	".", ";",
}

// A slice of available options for padding
var PaddingTypes = []string{PaddingTypeAdaptive, PaddingTypeFixed, PaddingTypeNone}

// A slice of available options for case transformation
var TransformTypes = []string{
	CaseTransformAlternate, CaseTransformAlternateLettercase, CaseTransformCapitalise,
	CaseTransformCapitaliseInvert, CaseTransformInvert, CaseTransformLower,
	CaseTransformLowerVowelUpperConsonant, CaseTransformNone, CaseTransformRandom,
	CaseTransformSentence, CaseTransformUpper,
}

var PaddingCharacterOptions = append([]string{PaddingCharacterRandom}, DefaultSpecialCharacters...)

var SeparatorCharacterOptions = append([]string{SeparatorCharacterRandom}, DefaultSpecialCharacters...)

// A slice of available word lists
var WordLists = []string{
	WordListAll, WordListDoctorWho, WordListEN, WordListENSmall, WordListGameOfThrones,
	WordListHarryPotter, WordListMiddleEarth, WordListStarTrek, WordListStarWars,
}

var WordListDescriptionMap = map[string]string{
	WordListAll:           "A combination of all word lists (44900+ words)",
	WordListDoctorWho:     "A Doctor Who word list (11300+ words)",
	WordListEN:            "A list of English words (14900+ words)",
	WordListENSmall:       "A small list of English words (8600+ words)",
	WordListGameOfThrones: "A Game of Thrones word list (8200+ words)",
	WordListHarryPotter:   "A Harry Potter word list (12500+ words)",
	WordListMiddleEarth:   "A Middle Earth word list containing words from The Hobbit, Lord of the Rings, The Silmarillion, and more (15400+ words)",
	WordListPokemon:       "A Pokemon word list (9000+ words)",
	WordListStarTrek:      "A Star Trek word list (8000+ words)",
	WordListStarWars:      "A Star Wars word list (12000+ words)",
}

var PresetDescriptionMap = map[string]string{
	PresetAppleID:       "A preset respecting the many prerequisites Apple places on Apple ID passwords. The preset also limits itself to symbols found on the iOS letter and number keyboards (i.e. not the awkward to reach symbol keyboard)",
	PresetDefault:       "The default preset resulting in a password consisting of 3 random words of between 4 and 8 letters with alternating case separated by a random character, with two random digits before and after, and padded with two random characters front and back",
	PresetNTLM:          "A preset for 14 character Windows NTLMv1 password. WARNING - only use this preset if you have to, it is too short to be acceptably secure",
	PresetSecurityQ:     "A preset for creating fake answers to security questions",
	PresetWeb16:         "A preset for websites that insist passwords not be longer than 16 characters",
	PresetWeb16XKPasswd: "A preset for websites that insist passwords not be longer than 16 characters, the same as the one found on xkpasswd.net.",
	PresetWeb32:         "A preset for websites that allow passwords up to 32 characteres long",
	PresetWiFi:          "A preset for generating 63 character long WPA2 keys",
	PresetXKCD:          "A preset for generating passwords similar to the example in the original XKCD cartoon, but with a dash to separate the four randomly capitalised words, two digits and a random special characters.",
	PresetXKCDXKPasswd:  "A preset for generating passwords similar to the example in the original XKCD cartoon, but with a dash to separate the four random words, and the capitalisation randomised to add sufficient entropy to avoid warnings.",
}
