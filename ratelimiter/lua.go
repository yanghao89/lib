package ratelimiter

import (
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	// redis 中执行的 lua脚本判断 key 是否被应该被限屏
	limiterScript = redis.NewScript(`
	redis.replicate_commands()
	redis.log(redis.LOG_DEBUG, "------------ ratelimiter script begin ------------"

	redis.log(redis.LOG_DEBUG, "------------ ratelimiter script end ------------")
`)
	RedisRatelimiterCacheExpiration = time.Minute * 60
)
