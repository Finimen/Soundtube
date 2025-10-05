package scripts

import "strings"

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	if !strings.Contains(email, ".") {
		return false
	}

	return true
}
