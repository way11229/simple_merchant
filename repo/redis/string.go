package redis

import (
	"context"
	"errors"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/way11229/simple_merchant/domain"
)

func (r *RedisClient) SetEx(ctx context.Context, input *domain.SetExParams) error {
	resp := r.rdb.SetEx(ctx, input.Key, input.Value, input.Expiration)
	if err := resp.Err(); err != nil {
		log.Printf("rdb.SetEx error = %v, params = %v", err, input)
		return domain.ErrUnknown
	}

	return nil
}

func (r *RedisClient) Get(ctx context.Context, input *domain.GetParams) (string, error) {
	resp, err := r.rdb.Get(ctx, input.Key).Result()
	if errors.Is(err, redis.Nil) {
		return "", domain.ErrRecordNotFound
	}

	if err != nil {
		log.Printf("rdb.Get error = %v, params = %v", err, input)
		return "", domain.ErrUnknown
	}

	return resp, nil
}

func (r *RedisClient) Del(ctx context.Context, input *domain.DelParams) error {
	if err := r.rdb.Del(ctx, input.Keys...).Err(); err != nil {
		log.Printf("rdb.Del error = %v, params = %v", err, input.Keys)
		return domain.ErrUnknown
	}

	return nil
}
