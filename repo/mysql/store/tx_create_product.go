package store

import (
	"context"
	"log"

	"github.com/way11229/simple_merchant/domain"
	mysql_sqlc "github.com/way11229/simple_merchant/repo/mysql/sqlc"
	"github.com/way11229/simple_merchant/utils"
)

func (s *SqlStore) TxCreateProduct(ctx context.Context, input *domain.MysqlTxCreateProductParams) (uint32, error) {
	var productId uint32
	err := s.execTx(ctx, func(q *mysql_sqlc.Queries) error {
		createProductResp, err := q.CreateProduct(ctx, input.CreateProductParams)
		if err != nil {
			log.Printf("CreateProduct error = %v, params = %v", err, input.CreateProductParams)
			return domain.ErrUnknown
		}

		lastInsertId, err := createProductResp.LastInsertId()
		if err != nil {
			log.Printf("LastInsertId error = %v", err)
			return domain.ErrUnknown
		}

		productId, err = utils.ConvertInt64ToUint32(lastInsertId)
		if err != nil {
			log.Printf("ConvertInt64ToUint32 error = %v, params = %d", err, lastInsertId)
			return domain.ErrUnknown
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return productId, nil
}
