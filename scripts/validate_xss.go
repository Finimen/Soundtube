package scripts

import "strings"

func ValidateXSS(field string) bool {
	if strings.Contains(field, "<") && strings.Contains(field, ">") {
		return true
	}

	return false
}
