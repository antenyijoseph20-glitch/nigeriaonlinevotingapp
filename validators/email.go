package validators

import "net/mail"

// ValidateEmail checks whether an email address is valid.
func ValidateEmail(email string) (bool, string) {

	_, err := mail.ParseAddress(email)
	if err != nil {
		return false, "Please enter a valid email address."
	}

	return true, ""
}
