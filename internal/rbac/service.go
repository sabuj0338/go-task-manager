package rbac

import (
	auth_repo "github.com/sabuj0338/go-task-manager/internal/auth/repository"
	"github.com/sabuj0338/go-task-manager/internal/models"
)

func CreateRole(dto CreateRoleDTO) error {
	role := &models.Role{
		Name: dto.Name,
	}
	return auth_repo.CreateRole(role)
}

func CreatePermission(dto CreatePermissionDTO) error {
	permission := &models.Permission{
		Name: dto.Name,
		Code: dto.Code,
	}
	return auth_repo.CreatePermission(permission)
}

func AssignPermissionToRole(dto AssignPermissionToRoleDTO) error {
	return auth_repo.AssignPermissionToRole(dto.RoleID, dto.PermissionID)
}

func AssignRoleToUser(dto AssignRoleToUserDTO) error {
	err := auth_repo.AssignRoleToUser(dto.UserID, dto.RoleID)
	if err != nil {
		return err
	}
	// Invalidate the user's permission cache after their roles change
	// to ensure new permissions are applied on their next request.
	return auth_repo.ClearUserPermissionCache(dto.UserID)
}
