package auth

import (
	"errors"

	"github.com/pquerna/otp/totp"
	"github.com/sabuj0338/go-task-manager/internal/auth/repository"
	"github.com/sabuj0338/go-task-manager/internal/models"
	"github.com/sabuj0338/go-task-manager/internal/utils"
	"github.com/sabuj0338/go-task-manager/pkg/lock"
	"github.com/sabuj0338/go-task-manager/pkg/token"
)

func Register(dto RegisterDTO) error {
	existingUser, _ := repository.GetUserByEmail(dto.Email)
	if existingUser != nil {
		return errors.New("user already exists")
	}

	hashed, _ := utils.HashPassword(dto.Password)

	user := &models.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: hashed,
		Role:     "user",
		Verified: false,
	}

	return repository.CreateUser(user)
}

func Login(dto LoginDTO) (*models.User, string, string, error) {
	user, err := repository.GetUserByEmail(dto.Email)
	if err != nil || user == nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	if lock.IsLocked(user.Email) {
		return nil, "", "", errors.New("Account temporarily locked. Try later")
	}

	if !utils.CheckPasswordHash(dto.Password, user.Password) {
		lock.LoginFailed(user.Email)
		return nil, "", "", errors.New("invalid credentials")
	}

	accessToken, err := token.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := token.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func GenerateTOTPSecret(email string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "GoTaskManager",
		AccountName: email,
	})
	if err != nil {
		return "", "", err
	}

	// key.Secret() → store this
	// key.URL() → use this to create QR

	return key.Secret(), key.URL(), nil
}

func VerifyTOTPToken(secret string, token string) bool {
	return totp.Validate(token, secret)
}

// func Login(dto LoginDTO) (string, error) {
// 	user, err := repository.GetUserByEmail(dto.Email)
// 	if err != nil || user == nil {
// 		return "", errors.New("invalid credentials")
// 	}

// 	if !utils.CheckPasswordHash(dto.Password, user.Password) {
// 		return "", errors.New("invalid credentials")
// 	}

// 	jwt, err := token.GenerateJWT(user.ID, user.Role)
// 	if err != nil {
// 		return "", err
// 	}

// 	return jwt, nil
// }
