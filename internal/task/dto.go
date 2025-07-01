package task

type CreateTaskDTO struct {
	Title       string `json:"title" validate:"required,min=3"`
	Description string `json:"description"`
}

type UpdateTaskDTO struct {
	Title       string `json:"title" validate:"omitempty,min=3"`
	Description string `json:"description"`
	Completed   *bool  `json:"completed"`
}
