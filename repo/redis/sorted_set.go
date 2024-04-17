package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/way11229/simple_merchant/domain"
)

func (r *RedisClient) ZAdd(ctx context.Context, input *domain.ZAddParams) error {
	members := []redis.Z{}
	for _, e := range input.Members {
		members = append(members, redis.Z{
			Score:  e.Score,
			Member: e.Member,
		})
	}

	if err := r.rdb.ZAdd(ctx, input.Key, members...); err != nil {
		log.Printf("ZAdd error = %v, params = %v", err, input)
		return domain.ErrUnknown
	}

	return nil
}

func (r *RedisClient) ZRem(ctx context.Context, input *domain.ZRemParams) error {
	if err := r.rdb.ZRem(ctx, input.Key, input.Members...); err != nil {
		log.Printf("ZRem error = %v, params = %v", err, input)
		return domain.ErrUnknown
	}

	return nil
}

func (r *RedisClient) ZRevRange(ctx context.Context, input *domain.ZRevRangeParams) ([]string, error) {
	resp, err := r.rdb.ZRevRange(ctx, input.Key, input.Start, input.Stop).Result()
	if err != nil {
		log.Printf("rdb.ZRevRange error = %v, params = %v", err, input)
		return nil, domain.ErrUnknown
	}

	return resp, nil
}
