package mfa

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

func SendEmailCode(email string) (string, error) {
	code := generateCode()
	key := "mfa:email:" + email
	err := redis.Client.Set(redis.Ctx, key, code, time.Minute*5).Err()
	if err != nil {
		return "", err
	}
	// TODO: integrate actual email sender here
	// fmt.Println("Sent email code to", email, ":", code)

	subject := "Verify your email"
	body := fmt.Sprintf(`<p>Your email verification code is: <b>%s</b></p>`, code)

	mail.Send(email, subject, body)

	return code, nil
}

func SendSMSCode(phone string) (string, error) {
	code := generateCode()
	key := "mfa:sms:" + phone
	err := redis.Client.Set(redis.Ctx, key, code, time.Minute*5).Err()
	if err != nil {
		return "", err
	}
	// TODO: integrate Twilio or other SMS API
	fmt.Println("Sent SMS code to", phone, ":", code)
	return code, nil
}

func VerifyCode(key string, code string) bool {
	val, err := redis.Client.Get(redis.Ctx, key).Result()
	if err != nil || val != code {
		return false
	}
	// Invalidate after use
	redis.Client.Del(redis.Ctx, key)
	return true
}
