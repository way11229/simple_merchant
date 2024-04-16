package domain

import (
	"context"
)

type CreateUserParams struct {
	Name     string
	Email    string
	Password string
}

func (c *CreateUserParams) Validate() error {
	if len([]rune(c.Name)) > USER_NAME_MAX_LEN {
		return ErrInvalidUserName
	}

	if len(c.Email) == 0 || !ValidateEmail(c.Email) {
		return ErrInvalidEmail
	}

	if len(c.Password) == 0 || !ValidateUserPassword(c.Password) {
		return ErrInvalidUserPassword
	}

	return nil
}

type CreateUserResult struct {
	UserId uint32
}

type UserService interface {
	CreateUser(ctx context.Context, input *CreateUserParams) (*CreateUserResult, error)
	DeleteUserById(ctx context.Context, userId uint32) error
}
