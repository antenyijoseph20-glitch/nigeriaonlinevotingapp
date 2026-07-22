package validators

import "unicode"

// ValidatePassword checks whether a password meets the required policy.
func ValidatePassword(password string) (bool, string) {

	if len(password) < 8 {
		return false, "Password must be at least 8 characters long."
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, ch := range password {

		switch {
		case unicode.IsUpper(ch):
			hasUpper = true

		case unicode.IsLower(ch):
			hasLower = true

		case unicode.IsDigit(ch):
			hasNumber = true

		default:
			hasSpecial = true
		}
	}

	if !hasUpper {
		return false, "Password must contain at least one uppercase letter."
	}

	if !hasLower {
		return false, "Password must contain at least one lowercase letter."
	}

	if !hasNumber {
		return false, "Password must contain at least one number."
	}

	if !hasSpecial {
		return false, "Password must contain at least one special character."
	}

	return true, ""
}
