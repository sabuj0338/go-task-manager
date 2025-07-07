package user

type CreateUserDTO struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,strong_password"`
}

type UpdateUserDTO struct {
	Name  string `json:"name" validate:"omitempty,min=3"`
	Email string `json:"email" validate:"omitempty,email"`
}
