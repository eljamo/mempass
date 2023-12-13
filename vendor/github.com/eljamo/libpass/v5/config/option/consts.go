package option

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
