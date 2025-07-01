package auth

type RegisterDTO struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
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
	Method string `json:"method" validate:"required,oneof=email sms totp"`
	Code   string `json:"code" validate:"required,len=6"`
	Trust  bool   `json:"trust"`
}

type ForgotPasswordDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordDTO struct {
	Email       string `json:"email" validate:"required,email"`
	Code        string `json:"code" validate:"required,len=6"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type EmailVerifyDTO struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}
