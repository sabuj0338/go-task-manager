package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sabuj0338/go-task-manager/internal/auth/repository"
	"github.com/sabuj0338/go-task-manager/pkg/response"
	"github.com/sabuj0338/go-task-manager/pkg/token"
)

// JWTProtected is a middleware to protect routes that require authentication.
func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Missing or malformed JWT")
		}

		const prefix = "Bearer "
		if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid authorization format")
		}
		tokenString := authHeader[len(prefix):]

		claims, err := token.ParseToken(tokenString)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid or expired JWT")
		}

		if claims["type"] != "access" {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid token type")
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid user ID in token")
		}

		c.Locals("user_id", uint(userID))
		return c.Next()
	}
}

// RequireRole is the old role-based middleware.
// DEPRECATED: Use RequirePermission for more granular control.
// This middleware is inefficient as it queries the database on every request without caching.
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(uint)
		if !ok {
			return response.Error(c, fiber.StatusForbidden, "Permission denied")
		}

		user, err := repository.GetUserByID(int(userID))
		if err != nil {
			return response.Error(c, fiber.StatusNotFound, "User not found")
		}

		hasRole := false
		for _, role := range user.Roles {
			if role.Name == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return response.Error(c, fiber.StatusForbidden, "Permission denied: requires '"+requiredRole+"' role")
		}

		return c.Next()
	}
}

// RequirePermission creates a middleware that checks if the authenticated user has a specific permission.
func RequirePermission(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(uint)
		if !ok {
			// This should not be reached if JWTProtected runs first, but it's a good safeguard.
			return response.Error(c, fiber.StatusForbidden, "Permission denied: user not authenticated")
		}

		// Get user permissions from cache or DB
		permissions, err := repository.GetUserPermissionsWithCache(userID)
		if err != nil {
			// It's good practice to log the actual error for debugging purposes.
			// log.Printf("Error getting user permissions for user %d: %v", userID, err)
			return response.Error(c, fiber.StatusInternalServerError, "Could not verify permissions")
		}
		fmt.Println("permissions:", permissions, userID)
		hasPermission := false
		for _, p := range permissions {
			if p == requiredPermission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return response.Error(c, fiber.StatusForbidden, "Permission denied: requires '"+requiredPermission+"' permission")
		}

		return c.Next()
	}
}
