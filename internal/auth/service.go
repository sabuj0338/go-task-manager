package auth

import (
	"errors"

	"github.com/pquerna/otp/totp"
	"github.com/sabuj0338/go-task-manager/internal/auth/repository"
	"github.com/sabuj0338/go-task-manager/internal/models"
	"github.com/sabuj0338/go-task-manager/internal/utils"
	"github.com/sabuj0338/go-task-manager/pkg/database"
	"github.com/sabuj0338/go-task-manager/pkg/lock"
	"github.com/sabuj0338/go-task-manager/pkg/token"
)

func Register(dto RegisterDTO) error {
	existingUser, _ := repository.GetUserByEmail(dto.Email)
	if existingUser != nil {
		return errors.New("user already exists")
	}

	// Find the default "user" role
	userRole, err := repository.GetRoleByName("user")
	if err != nil {
		return errors.New("default user role not found")
	}

	hashed, _ := utils.HashPassword(dto.Password)

	user := &models.User{
		Name:          dto.Name,
		Email:         dto.Email,
		Password:      hashed,
		EmailVerified: false,
	}

	if err := repository.CreateUser(user); err != nil {
		return err
	}

	return repository.AssignRoleToUser(user.ID, userRole.ID)
}

func Login(dto LoginDTO) (*models.User, string, string, error) {
	user, err := repository.GetUserByEmail(dto.Email)
	if err != nil || user == nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	if lock.IsLocked(user.Email) {
		return nil, "", "", errors.New("account temporarily locked. Try later")
	}

	if !utils.CheckPasswordHash(dto.Password, user.Password) {
		lock.LoginFailed(user.Email)
		return nil, "", "", errors.New("invalid credentials")
	}

	// accessToken, err := token.GenerateJWT(user.ID, user.Role)
	accessToken, err := token.GenerateJWT(user.ID)

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

func UpdateProfile(userID uint, dto UpdateProfileDTO) (*models.User, error) {
	emailExists := repository.CheckEmailExists(userID, dto.Email)
	if emailExists {
		return nil, errors.New("email already exists")
	}

	user, err := repository.GetUserByID(int(userID))
	if err != nil {
		return nil, err
	}

	if dto.Name != "" {
		user.Name = dto.Name
	}
	if dto.Phone != "" {
		user.Phone = dto.Phone
	}
	if dto.Email != "" {
		user.Email = dto.Email
	}
	if dto.CurrentPassword != "" && dto.NewPassword != "" {
		if !utils.CheckPasswordHash(dto.CurrentPassword, user.Password) {
			return nil, errors.New("invalid current password")
		}
		hashed, _ := utils.HashPassword(dto.NewPassword)
		user.Password = hashed
	}

	if err := database.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
