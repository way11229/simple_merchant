package store

import (
	"context"
	"log"

	"github.com/way11229/simple_merchant/domain"
	mysql_sqlc "github.com/way11229/simple_merchant/repo/mysql/sqlc"
	"github.com/way11229/simple_merchant/utils"
)

func (s *SqlStore) TxCreateUser(ctx context.Context, input *domain.MysqlTxCreateUserParams) (uint32, error) {
	var userId uint32
	err := s.execTx(ctx, func(q *mysql_sqlc.Queries) error {
		createUserResp, err := q.CreateUser(ctx, input.CreateUserParams)
		if err != nil {
			log.Printf("CreateUser error = %v, params = %v", err, input.CreateUserParams)
			return domain.ErrUnknown
		}

		lastInsertId, err := createUserResp.LastInsertId()
		if err != nil {
			log.Printf("LastInsertId error = %v", err)
			return domain.ErrUnknown
		}

		userId, err = utils.ConvertInt64ToUint32(lastInsertId)
		if err != nil {
			log.Printf("ConvertInt64ToUint32 error = %v, params = %d", err, lastInsertId)
			return domain.ErrUnknown
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return userId, nil
}
