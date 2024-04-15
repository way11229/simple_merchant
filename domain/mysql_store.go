package domain

import (
	mysql_sqlc "github.com/way11229/simple_merchant/repo/mysql/sqlc"
)

type MysqlStore interface {
	mysql_sqlc.Querier
}
