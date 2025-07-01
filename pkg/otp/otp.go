package otp

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

func SendResetCode(email string) (string, error) {
	code := generateCode()
	key := "reset:" + email
	err := redis.Client.Set(redis.Ctx, key, code, 5*time.Minute).Err()
	if err != nil {
		return "", err
	}
	// TODO: integrate real email sender
	// fmt.Println("Reset code sent to", email, ":", code)

	subject := "Forgot password email verification"
	body := fmt.Sprintf(`<p>Your email verification code is: <b>%s</b></p>`, code)

	mail.Send(email, subject, body)

	return code, nil
}

func VerifyResetCode(email, code string) bool {
	key := "reset:" + email
	val, err := redis.Client.Get(redis.Ctx, key).Result()
	if err != nil || val != code {
		return false
	}
	redis.Client.Del(redis.Ctx, key)
	return true
}
