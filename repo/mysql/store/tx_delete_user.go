package store

import (
	"context"
	"log"

	"github.com/way11229/simple_merchant/domain"
	mysql_sqlc "github.com/way11229/simple_merchant/repo/mysql/sqlc"
)

func (s *SqlStore) TxDeleteUser(ctx context.Context, input *domain.MysqlTxDeleteUserParams) error {
	return s.execTx(ctx, func(q *mysql_sqlc.Queries) error {
		if err := q.DeleteUserEmailVerificationCodeByUserId(ctx, input.UserId); err != nil {
			log.Printf("DeleteUserEmailVerificationCodeByUserId error = %v, params = %d", err, input.UserId)
			return domain.ErrUnknown
		}

		if err := q.DeleteUserAuthByUserId(ctx, input.UserId); err != nil {
			log.Printf("DeleteUserAuthByUserId error = %v, params = %d", err, input.UserId)
			return domain.ErrUnknown
		}

		if err := q.DeleteUserById(ctx, input.UserId); err != nil {
			log.Printf("DeleteUserById error = %v, params = %d", err, input.UserId)
			return domain.ErrUnknown
		}

		return nil
	})
}
