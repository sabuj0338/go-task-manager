package rbac

type CreateRoleDTO struct {
	Name string `json:"name" validate:"required"`
}

type CreatePermissionDTO struct {
	Name string `json:"name" validate:"required"`
	Code string `json:"code" validate:"required"`
}

type AssignPermissionToRoleDTO struct {
	RoleID       uint `json:"role_id" validate:"required"`
	PermissionID uint `json:"permission_id" validate:"required"`
}

type AssignRoleToUserDTO struct {
	UserID uint `json:"user_id" validate:"required"`
	RoleID uint `json:"role_id" validate:"required"`
}
