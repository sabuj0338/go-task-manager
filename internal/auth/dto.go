package auth

type RegisterDTO struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,strong_password"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type VerifyTOTPDTO struct {
	Token string `json:"token" validate:"required,len=6"`
}

type EnableMFAResponse struct {
	Secret  string `json:"secret"`
	QRImage string `json:"qr_image"`
}

type MFACodeVerifyDTO struct {
	VerificationToken string `json:"verification_token" validate:"required"`
	Method            string `json:"method" validate:"required,oneof=email sms totp"`
	Code              string `json:"code" validate:"required,len=6"`
	Trust             bool   `json:"trust"`
}

type DisableMFADTO struct {
	Password string `json:"password" validate:"required"`
}

type ForgotPasswordDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordDTO struct {
	Email       string `json:"email" validate:"required,email"`
	Code        string `json:"code" validate:"required,len=6"`
	NewPassword string `json:"new_password" validate:"required,strong_password"`
}

type ResetPasswordWithTokenDTO struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,strong_password"`
}

type EmailVerifyDTO struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

type UpdateProfileDTO struct {
	Name            string `json:"name" validate:"omitempty,min=3"`
	Phone           string `json:"phone" validate:"omitempty,min=1"`
	Email           string `json:"email" validate:"omitempty,email"`
	CurrentPassword string `json:"current_password" validate:"omitempty"`
	NewPassword     string `json:"new_password" validate:"omitempty,strong_password"`
}
