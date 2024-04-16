package domain

import (
	"context"

	mysql_sqlc "github.com/way11229/simple_merchant/repo/mysql/sqlc"
)

type MysqlTxCreateUserParams struct {
	CreateUserParams mysql_sqlc.CreateUserParams
}

type MysqlTxDeleteUserParams struct {
	UserId uint32
}

type MysqlTxCreateUserEmailVerificationCodeParams struct {
	CreateUserEmailVerificationCodeParams mysql_sqlc.CreateUserEmailVerificationCodeParams
	AfterCreate                           func() error
}

type MysqlStore interface {
	mysql_sqlc.Querier

	TxCreateUser(ctx context.Context, input *MysqlTxCreateUserParams) (uint32, error)
	TxDeleteUser(ctx context.Context, input *MysqlTxDeleteUserParams) error
	TxCreateUserEmailVerificationCode(ctx context.Context, input *MysqlTxCreateUserEmailVerificationCodeParams) error
}
