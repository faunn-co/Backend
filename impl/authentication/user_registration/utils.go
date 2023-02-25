package user_registration

import "unicode"

func isContainsNonNumeric(a string) bool {
	for _, char := range a {
		if !unicode.IsNumber(char) {
			return true
		}
	}
	return false
}
func isContainsSpecialChar(a string) bool {
	for _, char := range a {
		if unicode.IsSymbol(char) {
			return true
		}
	}
	for _, char := range a {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return true
		}
	}
	return false
}

func isContainsSpace(a string) bool {
	for _, char := range a {
		if unicode.IsSpace(char) {
			return true
		}
	}
	return false
}

func isContainsAtSign(a string) bool {
	match := "@"

	for _, char := range a {
		for _, key := range match {
			if char == key {
				return true
			}
		}
	}
	return false
}
