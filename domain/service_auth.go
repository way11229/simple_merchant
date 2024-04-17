package domain

import "context"

type LoginUserParams struct {
	Email    string
	Password string
}

func (l *LoginUserParams) Validate() error {
	if len(l.Email) == 0 || !ValidateEmail(l.Email) {
		return ErrInvalidEmail
	}

	if len(l.Password) == 0 {
		return ErrInvalidUserPassword
	}

	return nil
}

type LoginUserResult struct {
	Token            string
	EmailHasVerified bool
}

type CheckAccessTokenParams struct {
	AccessToken    string
	GetNowTimeFunc FuncTimeType
}

type CheckAccessTokenResult struct {
	UserId uint32
}

type LogoutUserParams struct {
	UserId uint32
}

type AuthService interface {
	LoginUser(ctx context.Context, input *LoginUserParams) (*LoginUserResult, error)
	CheckAccessToken(ctx context.Context, input *CheckAccessTokenParams) (*CheckAccessTokenResult, error)
	LogoutUser(ctx context.Context, input *LogoutUserParams) error
}
