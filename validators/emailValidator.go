package validators

import "regexp"

func IsEmailValid(email string) bool {

	// regular expression to match email addresses
	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

	// compile the regular expression
	pattern := regexp.MustCompile(regex)

	// match the email address against the regular expression
	return pattern.MatchString(email)
}