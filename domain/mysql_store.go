package domain

import (
	"context"

	mysql_sqlc "github.com/way11229/simple_merchant/repo/mysql/sqlc"
)

type MysqlTxCreateUserParams struct {
	CreateUserParams mysql_sqlc.CreateUserParams
}

type MysqlStore interface {
	mysql_sqlc.Querier

	TxCreateUser(ctx context.Context, input *MysqlTxCreateUserParams) (int64, error)
}
