package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sabuj0338/go-task-manager/internal/models"
	"github.com/sabuj0338/go-task-manager/pkg/database"
	"github.com/sabuj0338/go-task-manager/pkg/redis"
)

// GetRoleByName finds a role by its name.
func GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	if err := database.DB.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// AssignRoleToUser assigns a role to a user in the user_roles pivot table.
func AssignRoleToUser(userID uint, roleID uint) error {
	user := models.User{ID: userID}
	role := models.Role{ID: roleID}
	return database.DB.Model(&user).Association("Roles").Append(&role)
}

// GetUserRoles retrieves all roles for a given user.
func GetUserRoles(userID uint) ([]models.Role, error) {
	var roles []models.Role
	err := database.DB.
		Joins("JOIN user_roles on user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

// GetUserPermissionsWithCache retrieves all permission codes for a given user, with caching.
func GetUserPermissionsWithCache(userID uint) ([]string, error) {
	cacheKey := fmt.Sprintf("user_permissions:%d", userID)
	ctx := context.Background()

	// 1. Try to get from Redis cache
	cachedPermissions, err := redis.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		var permissions []string
		if json.Unmarshal([]byte(cachedPermissions), &permissions) == nil {
			return permissions, nil // Cache hit
		}
	}

	// 2. If cache miss, get from DB
	permissions, err := GetUserPermissions(userID)
	if err != nil {
		return nil, err
	}

	// 3. Store in Redis cache with a 15-minute expiration
	permissionsJSON, _ := json.Marshal(permissions)
	redis.Client.Set(ctx, cacheKey, permissionsJSON, 15*time.Minute)

	return permissions, nil
}

// GetUserPermissions retrieves all permission codes for a given user.
func GetUserPermissions(userID uint) ([]string, error) {
	var permissions []string
	err := database.DB.Table("permissions").
		Select("DISTINCT permissions.code").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Pluck("code", &permissions).Error

	return permissions, err
}

// CreateRole creates a new role.
func CreateRole(role *models.Role) error {
	return database.DB.Create(role).Error
}

// CreatePermission creates a new permission.
func CreatePermission(permission *models.Permission) error {
	return database.DB.Create(permission).Error
}

// AssignPermissionToRole assigns a permission to a role.
func AssignPermissionToRole(roleID uint, permissionID uint) error {
	role := models.Role{ID: roleID}
	permission := models.Permission{ID: permissionID}
	return database.DB.Model(&role).Association("Permissions").Append(&permission)
}

// GetRoleByID finds a role by its ID.
func GetRoleByID(id uint) (*models.Role, error) {
	var role models.Role
	err := database.DB.First(&role, id).Error
	return &role, err
}

// ClearUserPermissionCache removes a user's cached permissions from Redis.
func ClearUserPermissionCache(userID uint) error {
	cacheKey := fmt.Sprintf("user_permissions:%d", userID)
	return redis.Client.Del(context.Background(), cacheKey).Err()
}
