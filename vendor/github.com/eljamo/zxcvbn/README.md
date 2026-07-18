[![GoDoc](https://godoc.org/github.com/eljamo/zxcvbn?status.svg)](https://godoc.org/github.com/eljamo/zxcvbn)

# zxcvbn

This project is a fork of [trustelem/zxcvbn](https://github.com/trustelem/zxcvbn), the Go port of [dropbox/zxcvbn](https://github.com/dropbox/zxcvbn).

It estimates password strength by looking at how real-world password crackers work. Instead of relying only on length or character rules, it checks for common patterns such as frequently used passwords, names and surnames from U.S. Census data, popular English words, dates, repeated characters, sequences, keyboard walks like qwerty, and l33t substitutions. It then uses those matches to give a conservative estimate of how difficult the password would be to guess.

While [trustelem/zxcvbn](https://github.com/trustelem/zxcvbn) focused on being a 1:1 port of [dropbox/zxcvbn](https://github.com/dropbox/zxcvbn), including matching its output, this fork prioritises improvements over strict output parity.

# Modifications

- Improved year detection so recent years are handled dynamically and the regex no longer needs updating for each new decade.
- Improved l33t matching for mixed substitutions, so variants like `p@5$w0rd` are recognized as `password`

# Fork Comparison

## Basic Example of trustelem/zxcvbn

## Program

```go
package main

import (
	"fmt"

	"github.com/trustelem/zxcvbn"
)

func main() {
	password := "p@5$w0rd"

	result := zxcvbn.PasswordStrength(password, nil)

	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Score: %d/4\n", result.Score)
	fmt.Printf("Guesses: %.0f\n", result.Guesses)
	fmt.Printf("Calculation time: %.3fs\n", result.CalcTime)
}
```

## Output

```
~/c/g/zxcvbn-trustelem (master|✚1) $ go run ./example/basic/main.go
Password: p@5$w0rd
Score: 2/4
Guesses: 12330000
Calculation time: 0.000s
```

## Basic Example of eljamo/zxcvbn

## Program

```go
package main

import (
	"fmt"

	"github.com/eljamo/zxcvbn"
)

func main() {
	password := "p@5$w0rd"

	result := zxcvbn.PasswordStrength(password, nil)

	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Score: %d/4\n", result.Score)
	fmt.Printf("Guesses: %.0f\n", result.Guesses)
	fmt.Printf("Guesses (log10): %.2f\n", result.GuessesLog10)
	fmt.Printf("Online throttled crack time: %s\n", result.CrackTimesDisplay["online_throttling_100_per_hour"])
	fmt.Printf("Online unthrottled crack time: %s\n", result.CrackTimesDisplay["online_no_throttling_10_per_second"])
	fmt.Printf("Offline fast hash crack time: %s\n", result.CrackTimesDisplay["offline_fast_hashing_1e10_per_second"])
	fmt.Printf("Offline slow hash crack time: %s\n", result.CrackTimesDisplay["offline_slow_hashing_1e4_per_second"])
	fmt.Printf("Calculation time: %.3fs\n", result.CalcTime)
}
```

## Output

```
~/c/g/zxcvbn (main|✔) $ go run ./example/basic/main.go
Password: p@5$w0rd
Score: 0/4
Guesses: 33
Guesses (log10): 1.52
Online throttled crack time: 20 minutes
Online unthrottled crack time: 3 seconds
Offline fast hash crack time: less than a second
Offline slow hash crack time: less than a second
Feedback Warning: This is similar to a commonly used password.
Feedback Suggestions: [Predictable substitutions like '@' instead of 'a' don't help very much. Add another word or two. Uncommon words are better.]
Calculation time: 0.000s
```
