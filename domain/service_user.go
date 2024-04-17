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

type DeleteUserByIdParams struct {
	UserId uint32
}

type GetUserEmailVerificationCodeParams struct {
	Email string
}

func (g *GetUserEmailVerificationCodeParams) Validate() error {
	if len(g.Email) == 0 || !ValidateEmail(g.Email) {
		return ErrInvalidEmail
	}

	return nil
}

type VerifyUserEmailParams struct {
	Email                 string
	EmailVerificationCode string
	GetNowTimeFunc        FuncTimeType
}

func (g *VerifyUserEmailParams) Validate() error {
	if len(g.Email) == 0 || !ValidateEmail(g.Email) {
		return ErrInvalidEmail
	}

	if len(g.EmailVerificationCode) == 0 {
		return ErrInvalidVerificationCode
	}

	return nil
}

type UserService interface {
	CreateUser(ctx context.Context, input *CreateUserParams) (*CreateUserResult, error)
	DeleteUserById(ctx context.Context, input *DeleteUserByIdParams) error
	GetUserEmailVerificationCode(ctx context.Context, input *GetUserEmailVerificationCodeParams) error
	VerifyUserEmail(ctx context.Context, input *VerifyUserEmailParams) error
}
