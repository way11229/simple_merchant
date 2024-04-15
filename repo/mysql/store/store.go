package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/way11229/simple_merchant/domain"
	mysql_sqlc "github.com/way11229/simple_merchant/repo/mysql/sqlc"
)

type SqlStore struct {
	db *sql.DB
	*mysql_sqlc.Queries
}

func NewStore(db *sql.DB) domain.MysqlStore {
	return &SqlStore{
		db:      db,
		Queries: mysql_sqlc.New(db),
	}
}

func (s *SqlStore) execTx(ctx context.Context, fn func(*mysql_sqlc.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := mysql_sqlc.New(tx)
	if err := fn(q); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}
