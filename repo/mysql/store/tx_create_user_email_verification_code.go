package store

import (
	"context"
	"log"

	"github.com/way11229/simple_merchant/domain"
	mysql_sqlc "github.com/way11229/simple_merchant/repo/mysql/sqlc"
)

func (s *SqlStore) TxCreateUserEmailVerificationCode(ctx context.Context, input *domain.MysqlTxCreateUserEmailVerificationCodeParams) error {
	return s.execTx(ctx, func(q *mysql_sqlc.Queries) error {
		if err := q.CreateUserEmailVerificationCode(ctx, input.CreateUserEmailVerificationCodeParams); err != nil {
			log.Printf("CreateUserEmailVerificationCode error = %v, params = %v", err, input.CreateUserEmailVerificationCodeParams)
			return domain.ErrUnknown
		}

		return input.AfterCreate()
	})
}
