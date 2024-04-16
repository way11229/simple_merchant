package domain

import (
	"net/mail"
	"strings"
	"unicode"
)

const (
	USER_NAME_MAX_LEN                = 100
	USER_PASSWORD_MIN_LEN            = 6
	USER_PASSWORD_MAX_LEN            = 16
	USER_PASSWORD_SPECIAL_CHARACTERS = "()[]{}<>+-*/?,.:;\"'_\\|~`!@#$%^&="
)

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidateUserPassword(pwd string) bool {
	pwdLen := len(pwd)
	if pwdLen > USER_PASSWORD_MAX_LEN || pwdLen < USER_PASSWORD_MIN_LEN {
		return false
	}

	hasUpperChar := false
	hasLowerChar := false
	hasSpecialChar := false

	for _, r := range pwd {
		switch {
		case unicode.IsLower(r):
			hasLowerChar = true
		case unicode.IsUpper(r):
			hasUpperChar = true
		case strings.ContainsRune(USER_PASSWORD_SPECIAL_CHARACTERS, r):
			hasSpecialChar = true
		default:
		}
	}

	return hasUpperChar && hasLowerChar && hasSpecialChar
}
