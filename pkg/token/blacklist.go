package token

import (
	"time"

	"github.com/sabuj0338/go-task-manager/pkg/redis"
)

func BlacklistRefreshToken(token string, ttl time.Duration) error {
	key := "blacklist:" + token
	return redis.Client.Set(redis.Ctx, key, "1", ttl).Err()
}

func IsRefreshTokenBlacklisted(token string) bool {
	key := "blacklist:" + token
	val, err := redis.Client.Get(redis.Ctx, key).Result()
	return err == nil && val == "1"
}
