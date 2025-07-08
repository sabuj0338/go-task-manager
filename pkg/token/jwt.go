package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// func GenerateJWT(userID uint, role string) (string, error) {
// 	claims := jwt.MapClaims{
// 		"user_id": userID,
// 		"role":    role,
// 		"exp":     time.Now().Add(time.Minute * 60).Unix(),
// 	}

// 	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return t.SignedString([]byte(os.Getenv("JWT_SECRET")))
// }

func GenerateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		// "role":    role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Standard 24-hour access token
		"type": "access",
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// GenerateMFAVerificationToken creates a short-lived token for the MFA verification step.
func GenerateMFAVerificationToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 5).Unix(), // Short-lived: 5 minutes
		"type":    "mfa_verification",
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// GeneratePasswordResetToken creates a short-lived token for resetting password.
func GeneratePasswordResetToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 15).Unix(), // Short-lived: 15 minutes
		"type":    "password_reset",
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
