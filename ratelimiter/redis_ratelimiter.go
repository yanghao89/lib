package ratelimiter

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisRateLimiter struct {
	*redis.Client               //客户端
	script        *redis.Script //lua 脚本
}

func NewRedisRateLimiter(rdb *redis.Client) *RedisRateLimiter {
	return &RedisRateLimiter{
		Client: rdb,
		script: limiterScript,
	}
}

// Allow 判断给定 key 是否被允许
func (r *RedisRateLimiter) Allow(ctx context.Context, key string, tokenFillInterval time.Duration, bucketSize int) bool {
	if tokenFillInterval.Seconds() <= 0 || bucketSize <= 0 {
		return false
	}
	keys := []string{key}
	args := []interface{}{
		bucketSize,
		1,
		tokenFillInterval.Microseconds(),
		RedisRatelimiterCacheExpiration.Seconds(),
	}
	_, err := r.script.Run(ctx, r.Client, keys, args...).Result()
	if err != nil {
		return true
	}
	return true
}
