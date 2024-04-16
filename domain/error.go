package domain

import "errors"

var (
	ErrUnknown                           = errors.New("unknown error")
	ErrPermissionDeny                    = errors.New("permission deny")
	ErrMissingRequiredParameter          = errors.New("missing required parameters")
	ErrLoginAborted                      = errors.New("login aborted")
	ErrInvalidUserName                   = errors.New("invalid user name")
	ErrInvalidEmail                      = errors.New("invalid email")
	ErrInvalidUserPassword               = errors.New("invalid user password")
	ErrUserEmailDuplicated               = errors.New("user email duplicated")
	ErrEmailHasVerified                  = errors.New("email has verified")
	ErrInvalidVerificationCode           = errors.New("invalid verification code")
	ErrVerificationCodeExpired           = errors.New("verification code expired")
	ErrSendVerificationCodeInShortPeriod = errors.New("send verification code in short period")
	ErrRecordNotFound                    = errors.New("record not found")
)
