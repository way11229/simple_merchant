package domain

import (
	"context"
	"time"
)

type SetExpirationParams struct {
	Key        string
	Expiration time.Duration
}

type ZAddParams struct {
	Key     string
	Members []*ZScoreMember
}

type ZScoreMember struct {
	Score  float64
	Member interface{}
}

type ZRemParams struct {
	Key     string
	Members []interface{}
}

type ZRevRangeParams struct {
	Key   string
	Start int64
	Stop  int64
}

type SetExParams struct {
	Key        string
	Value      interface{}
	Expiration time.Duration
}

type GetParams struct {
	Key string
}

type RedisClient interface {
	SetExpiration(ctx context.Context, input *SetExpirationParams) error

	// sorted set
	ZAdd(ctx context.Context, input *ZAddParams) error
	ZRem(ctx context.Context, input *ZRemParams) error
	ZRevRange(ctx context.Context, input *ZRevRangeParams) ([]string, error)

	// string
	SetEx(ctx context.Context, input *SetExParams) error
	Get(ctx context.Context, input *GetParams) (string, error)
}
