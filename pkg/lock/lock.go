package lock

import (
	"fmt"
	"time"

	"github.com/sabuj0338/go-task-manager/pkg/redis"
)

func LoginFailed(email string) {
	key := "fail:login:" + email
	redis.Client.Incr(redis.Ctx, key)
	redis.Client.Expire(redis.Ctx, key, 10*time.Minute)
}

func MFAFailed(userID uint) {
	key := fmt.Sprintf("fail:mfa:%d", userID)
	redis.Client.Incr(redis.Ctx, key)
	redis.Client.Expire(redis.Ctx, key, 10*time.Minute)
}

func IsLocked(email string) bool {
	key := "fail:login:" + email
	countStr, _ := redis.Client.Get(redis.Ctx, key).Result()
	return countStr >= "5" // lock after 5 failed attempts
}

func LockoutRemaining(email string) int {
	key := "fail:login:" + email
	val, _ := redis.Client.Get(redis.Ctx, key).Int()
	return 5 - val
}
