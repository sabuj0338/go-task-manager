package verify

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sabuj0338/go-task-manager/pkg/mail"
	"github.com/sabuj0338/go-task-manager/pkg/redis"
)

func generateCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func SendVerificationEmail(email string) error {
	code := generateCode()
	key := "verify:" + email
	err := redis.Client.Set(redis.Ctx, key, code, 5*time.Minute).Err()
	if err != nil {
		return err
	}

	subject := "Verify your email"
	body := fmt.Sprintf(`<p>Your email verification code is: <b>%s</b></p>`, code)

	return mail.Send(email, subject, body)
}

func VerifyCode(email, code string) bool {
	key := "verify:" + email
	val, err := redis.Client.Get(redis.Ctx, key).Result()
	if err != nil || val != code {
		return false
	}
	redis.Client.Del(redis.Ctx, key)
	return true
}
