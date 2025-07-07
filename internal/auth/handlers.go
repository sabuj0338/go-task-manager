package auth

import (
	"strconv"
	"time"
	"unicode"

	"github.com/sabuj0338/go-task-manager/internal/auth/repository"
	"github.com/sabuj0338/go-task-manager/pkg/database"
	"github.com/sabuj0338/go-task-manager/pkg/lock"
	"github.com/sabuj0338/go-task-manager/pkg/mfa"
	"github.com/sabuj0338/go-task-manager/pkg/otp"
	"github.com/sabuj0338/go-task-manager/pkg/response"
	"github.com/sabuj0338/go-task-manager/pkg/token"
	"github.com/sabuj0338/go-task-manager/pkg/verify"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("strong_password", strongPasswordValidation)
}

func strongPasswordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}

func RegisterHandler(c *fiber.Ctx) error {
	var dto RegisterDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	if err := validate.Struct(dto); err != nil {
		return response.ValidationErrorResponse(c, err)
	}

	if err := Register(dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Registration successful", nil)
}

func LoginHandler(c *fiber.Ctx) error {
	var dto LoginDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	if err := validate.Struct(dto); err != nil {
		return response.ValidationErrorResponse(c, err)
	}

	user, accessToken, refreshToken, err := Login(dto)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	// Check if MFA is required
	_, enabled, _ := repository.GetUserMFA(user.ID)

	// Optional: if trusted, skip MFA
	trusted := c.Cookies("trusted_device_"+strconv.Itoa(int(user.ID))) == "1"

	if enabled && !trusted {
		return response.Success(c, fiber.StatusOK, "MFA required", fiber.Map{"mfa": "totp_required"})
	}

	if !enabled && !trusted {
		// fallback to email code
		_, err := mfa.SendEmailCode(user.Email)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "Failed to send MFA code")
		}
		return response.Success(c, fiber.StatusOK, "MFA required", fiber.Map{"mfa": "email_code_required", "user_id": user.ID})
	}

	return response.Success(c, fiber.StatusOK, "Logged in successfully", fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

// func RefreshTokenHandler(c *fiber.Ctx) error {
// 	type Body struct {
// 		RefreshToken string `json:"refresh_token" validate:"required"`
// 	}
// 	var body Body
// 	if err := c.BodyParser(&body); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
// 	}
// 	if err := validate.Struct(body); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	parsedToken, err := token.VerifyRefreshToken(body.RefreshToken)
// 	if err != nil || !parsedToken.Valid {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
// 	}

// 	claims := parsedToken.Claims.(jwt.MapClaims)
// 	if claims["type"] != "refresh" {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token type"})
// 	}

// 	userID := uint(claims["user_id"].(float64))
// 	// userEmail := claims["email"].(string)
// 	newAccessToken, err := token.GenerateJWT(userID, "user") // Ideally fetch role again
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate access token"})
// 	}

// 	return c.JSON(fiber.Map{
// 		"access_token": newAccessToken,
// 	})
// }

func RefreshTokenHandler(c *fiber.Ctx) error {
	type Payload struct {
		RefreshToken string `json:"refresh_token"`
	}
	var p Payload
	if err := c.BodyParser(&p); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Missing refresh token")
	}

	if token.IsRefreshTokenBlacklisted(p.RefreshToken) {
		return response.Error(c, fiber.StatusForbidden, "Token is blacklisted")
	}

	claims, err := token.ParseToken(p.RefreshToken)
	if err != nil {
		return response.Error(c, fiber.StatusForbidden, "Invalid refresh token")
	}

	userID := uint(claims["user_id"].(float64))
	role := claims["role"].(string)

	newAccessToken, _ := token.GenerateJWT(userID, role)
	newRefreshToken, _ := token.GenerateRefreshToken(userID)

	return response.Success(c, fiber.StatusOK, "Token refreshed successfully", fiber.Map{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func SetupMFAHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	// email := c.Locals("email").(string)
	user, err := repository.GetUserByID(int(userID))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch user")
	}

	secret, qrURL, err := GenerateTOTPSecret(user.Email)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to generate secret")
	}

	err = repository.UpdateUserMFASecret(userID, secret)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to store MFA")
	}

	return response.Success(c, fiber.StatusOK, "MFA setup successful", fiber.Map{
		"secret":   secret,
		"qr_image": qrURL,
	})
}

func VerifyMFAHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var dto VerifyTOTPDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid input")
	}
	if err := validate.Struct(dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	secret, enabled, err := repository.GetUserMFA(userID)
	if err != nil || !enabled {
		return response.Error(c, fiber.StatusBadRequest, "MFA not enabled")
	}

	if !VerifyTOTPToken(secret, dto.Token) {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid token")
	}

	return response.Success(c, fiber.StatusOK, "MFA verified", nil)
}

func VerifyMFACodeHandler(c *fiber.Ctx) error {
	var dto MFACodeVerifyDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid input")
	}
	if err := validate.Struct(dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	userIDStr := c.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Missing user_id")
	}

	user, _ := repository.GetUserByID(userID)

	if dto.Method != "totp" {
		var key string
		if dto.Method == "email" {
			key = "mfa:email:" + user.Email
		} else {
			key = "mfa:sms:" + user.Phone
		}

		// if !mfa.VerifyCode(key, dto.Code) {
		// 	return c.Status(401).JSON(fiber.Map{"error": "Invalid MFA code"})
		// }
		if !mfa.VerifyCode(key, dto.Code) {
			lock.MFAFailed(uint(userID))
			return response.Error(c, fiber.StatusForbidden, "Invalid MFA code")
		}
	} else {
		secret, enabled, err := repository.GetUserMFA(uint(userID))
		if err != nil || !enabled {
			return response.Error(c, fiber.StatusBadRequest, "MFA not enabled")
		}
		if !VerifyTOTPToken(secret, dto.Code) {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid token")
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

	return response.Success(c, fiber.StatusOK, "MFA verified", fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

func ForgotPasswordHandler(c *fiber.Ctx) error {
	var dto ForgotPasswordDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid input")
	}
	if err := validate.Struct(dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	user, err := repository.GetUserByEmail(dto.Email)
	if err != nil || user == nil {
		return response.Error(c, fiber.StatusNotFound, "User not found")
	}

	_, err = otp.SendResetCode(dto.Email)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to send reset code")
	}

	return response.Success(c, fiber.StatusOK, "Reset code sent to your email", nil)
}

func ResetPasswordHandler(c *fiber.Ctx) error {
	var dto ResetPasswordDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid input")
	}
	if err := validate.Struct(dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	ok := otp.VerifyResetCode(dto.Email, dto.Code)
	if !ok {
		return response.Error(c, fiber.StatusForbidden, "Invalid or expired code")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(dto.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Password hashing failed")
	}

	query := `UPDATE users SET password = ?, updated_at = NOW() WHERE email = ?`
	_, err = database.DB.Exec(query, string(hashed), dto.Email)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to reset password")
	}

	return response.Success(c, fiber.StatusOK, "Password reset successful", nil)
}

func SendEmailVerificationHandler(c *fiber.Ctx) error {
	var dto ForgotPasswordDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid input")
	}
	if err := validate.Struct(dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	user, err := repository.GetUserByEmail(dto.Email)
	if err != nil || user == nil {
		return response.Error(c, fiber.StatusNotFound, "User not found")
	}

	err = verify.SendVerificationEmail(dto.Email)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to send email")
	}

	return response.Success(c, fiber.StatusOK, "Verification code sent", nil)
}

func VerifyEmailCodeHandler(c *fiber.Ctx) error {
	var dto EmailVerifyDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid input")
	}
	if err := validate.Struct(dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	if !verify.VerifyCode(dto.Email, dto.Code) {
		return response.Error(c, fiber.StatusForbidden, "Invalid or expired code")
	}

	// mark verified
	query := `UPDATE users SET email_verified = true, updated_at = NOW() WHERE email = ?`
	_, err := database.DB.Exec(query, dto.Email)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update verification status")
	}

	return response.Success(c, fiber.StatusOK, "Email verified successfully", nil)
}

func LogoutHandler(c *fiber.Ctx) error {
	type Payload struct {
		RefreshToken string `json:"refresh_token"`
	}
	var p Payload
	if err := c.BodyParser(&p); err != nil || p.RefreshToken == "" {
		return response.Error(c, fiber.StatusBadRequest, "Missing refresh token")
	}

	claims, err := token.ParseToken(p.RefreshToken)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid token")
	}

	exp := time.Unix(int64(claims["exp"].(float64)), 0)
	ttl := time.Until(exp)

	err = token.BlacklistRefreshToken(p.RefreshToken, ttl)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to blacklist token")
	}

	return response.Success(c, fiber.StatusOK, "Logged out successfully", nil)
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
