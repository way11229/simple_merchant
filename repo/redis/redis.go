package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/way11229/simple_merchant/domain"
)

type RedisClient struct {
	rdb *redis.Client
}

func NewRedisClient(addr, pwd string) domain.RedisClient {
	return &RedisClient{
		rdb: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pwd,
			DB:       0,
		}),
	}
}

func (r *RedisClient) SetExpiration(ctx context.Context, input *domain.SetExpirationParams) error {
	resp, err := r.rdb.Expire(ctx, input.Key, input.Expiration).Result()
	if err != nil {
		log.Printf("rdb.Expire error = %v, params = %v", err, input)
		return domain.ErrUnknown
	}

	if !resp {
		log.Print("redis set expiration fail")
		return domain.ErrUnknown
	}

	return nil
}
