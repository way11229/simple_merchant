// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package sqlc

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateProduct(ctx context.Context, arg CreateProductParams) (sql.Result, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error)
	CreateUserAuth(ctx context.Context, arg CreateUserAuthParams) error
	CreateUserEmailVerificationCode(ctx context.Context, arg CreateUserEmailVerificationCodeParams) error
	DecreaseUserEmailVerificationCodeMaxTryById(ctx context.Context, id uint32) error
	DeleteUserAuthByUserId(ctx context.Context, userID uint32) error
	DeleteUserById(ctx context.Context, id uint32) error
	DeleteUserEmailVerificationCodeByUserId(ctx context.Context, userID uint32) error
	GetLastCreatedUserEmailVerificationCodeByUserId(ctx context.Context, userID uint32) (UserEmailVerificationCode, error)
	GetProductById(ctx context.Context, id uint32) (Product, error)
	GetUserAuthByUserId(ctx context.Context, userID uint32) (UserAuth, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id uint32) (User, error)
	GetUserEmailVerificationCodeByEmailAndVerificationCode(ctx context.Context, arg GetUserEmailVerificationCodeByEmailAndVerificationCodeParams) (UserEmailVerificationCode, error)
	ListTheRecommendedProducts(ctx context.Context, arg ListTheRecommendedProductsParams) ([]Product, error)
	UpdateUserAuthById(ctx context.Context, arg UpdateUserAuthByIdParams) error
	VerifyUserEmailById(ctx context.Context, id uint32) error
}

var _ Querier = (*Queries)(nil)
