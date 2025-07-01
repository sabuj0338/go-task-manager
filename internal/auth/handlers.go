package auth

import (
	"strconv"
	"time"

	"github.com/sabuj0338/go-task-manager/internal/auth/repository"
	"github.com/sabuj0338/go-task-manager/pkg/database"
	"github.com/sabuj0338/go-task-manager/pkg/mfa"
	"github.com/sabuj0338/go-task-manager/pkg/otp"
	"github.com/sabuj0338/go-task-manager/pkg/token"
	"github.com/sabuj0338/go-task-manager/pkg/verify"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var validate = validator.New()

func RegisterHandler(c *fiber.Ctx) error {
	var dto RegisterDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := validate.Struct(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := Register(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Registration successful"})
}

func LoginHandler(c *fiber.Ctx) error {
	var dto LoginDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := validate.Struct(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, accessToken, refreshToken, err := Login(dto)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if MFA is required
	_, enabled, _ := repository.GetUserMFA(user.ID)

	// Optional: if trusted, skip MFA
	trusted := c.Cookies("trusted_device_"+strconv.Itoa(int(user.ID))) == "1"

	if enabled && !trusted {
		return c.JSON(fiber.Map{"mfa": "totp_required"})
	}

	if !enabled && !trusted {
		// fallback to email code
		_, err := mfa.SendEmailCode(user.Email)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to send MFA code"})
		}
		return c.JSON(fiber.Map{"mfa": "email_code_required", "user_id": user.ID})
	}

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

func RefreshTokenHandler(c *fiber.Ctx) error {
	type Body struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	var body Body
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	parsedToken, err := token.VerifyRefreshToken(body.RefreshToken)
	if err != nil || !parsedToken.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	claims := parsedToken.Claims.(jwt.MapClaims)
	if claims["type"] != "refresh" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token type"})
	}

	userID := uint(claims["user_id"].(float64))
	// userEmail := claims["email"].(string)
	newAccessToken, err := token.GenerateJWT(userID, "user") // Ideally fetch role again
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate access token"})
	}

	return c.JSON(fiber.Map{
		"access_token": newAccessToken,
	})
}

func SetupMFAHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	// email := c.Locals("email").(string)
	user, err := repository.GetUserByID(int(userID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch user"})
	}

	secret, qrURL, err := GenerateTOTPSecret(user.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	err = repository.UpdateUserMFASecret(userID, secret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to store MFA"})
	}

	return c.JSON(fiber.Map{
		"secret":   secret,
		"qr_image": qrURL,
	})
}

func VerifyMFAHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var dto VerifyTOTPDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}
	if err := validate.Struct(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	secret, enabled, err := repository.GetUserMFA(userID)
	if err != nil || !enabled {
		return c.Status(400).JSON(fiber.Map{"error": "MFA not enabled"})
	}

	if !VerifyTOTPToken(secret, dto.Token) {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	return c.JSON(fiber.Map{"message": "MFA verified"})
}

func VerifyMFACodeHandler(c *fiber.Ctx) error {
	var dto MFACodeVerifyDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if err := validate.Struct(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	userIDStr := c.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Missing user_id"})
	}

	user, _ := repository.GetUserByID(userID)

	if dto.Method != "totp" {
		var key string
		if dto.Method == "email" {
			key = "mfa:email:" + user.Email
		} else {
			key = "mfa:sms:" + user.Phone
		}

		if !mfa.VerifyCode(key, dto.Code) {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid MFA code"})
		}
	} else {
		secret, enabled, err := repository.GetUserMFA(uint(userID))
		if err != nil || !enabled {
			return c.Status(400).JSON(fiber.Map{"error": "MFA not enabled"})
		}
		if !VerifyTOTPToken(secret, dto.Code) {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid token"})
		}
	}

	// MFA passed → generate token
	accessToken, _ := token.GenerateJWT(uint(userID), "user")
	refreshToken, _ := token.GenerateRefreshToken(uint(userID))

	// Optional: trust device (set cookie)
	if dto.Trust {
		c.Cookie(&fiber.Cookie{
			Name:     "trusted_device_" + strconv.Itoa(userID),
			Value:    "1",
			Expires:  time.Now().Add(7 * 24 * time.Hour),
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Strict",
		})
	}

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

func ForgotPasswordHandler(c *fiber.Ctx) error {
	var dto ForgotPasswordDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}
	if err := validate.Struct(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := repository.GetUserByEmail(dto.Email)
	if err != nil || user == nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

	_, err = otp.SendResetCode(dto.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send reset code"})
	}

	return c.JSON(fiber.Map{"message": "Reset code sent to your email"})
}

func ResetPasswordHandler(c *fiber.Ctx) error {
	var dto ResetPasswordDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}
	if err := validate.Struct(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	ok := otp.VerifyResetCode(dto.Email, dto.Code)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired code"})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(dto.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Password hashing failed"})
	}

	query := `UPDATE users SET password = ?, updated_at = NOW() WHERE email = ?`
	_, err = database.DB.Exec(query, string(hashed), dto.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to reset password"})
	}

	return c.JSON(fiber.Map{"message": "Password reset successful"})
}

func SendEmailVerificationHandler(c *fiber.Ctx) error {
	var dto ForgotPasswordDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if err := validate.Struct(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := repository.GetUserByEmail(dto.Email)
	if err != nil || user == nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	err = verify.SendVerificationEmail(dto.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send email"})
	}

	return c.JSON(fiber.Map{"message": "Verification code sent"})
}

func VerifyEmailCodeHandler(c *fiber.Ctx) error {
	var dto EmailVerifyDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if err := validate.Struct(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if !verify.VerifyCode(dto.Email, dto.Code) {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired code"})
	}

	// mark verified
	query := `UPDATE users SET email_verified = true, updated_at = NOW() WHERE email = ?`
	_, err := database.DB.Exec(query, dto.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update verification status"})
	}

	return c.JSON(fiber.Map{"message": "Email verified successfully"})
}

// func LoginHandler(c *fiber.Ctx) error {
// 	var dto LoginDTO
// 	if err := c.BodyParser(&dto); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
// 	}

// 	if err := validate.Struct(dto); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	token, err := Login(dto)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	return c.JSON(fiber.Map{"token": token})
// }
