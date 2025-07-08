package middleware

// Requires Go 1.21+

// // RequirePermission checks if the authenticated user has a specific permission.
// // It is highly efficient due to caching of user permissions.
// func RequirePermission(requiredPermission string) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		userID, ok := c.Locals("user_id").(uint)
// 		if !ok {
// 			return response.Error(c, fiber.StatusForbidden, "Permission denied: Invalid user context")
// 		}

// 		// Get user permissions from cache or database
// 		permissions, err := auth_repo.GetUserPermissionsWithCache(userID)
// 		if err != nil {
// 			return response.Error(c, fiber.StatusInternalServerError, "Could not retrieve user permissions")
// 		}

// 		// Check if the user has the required permission
// 		if !slices.Contains(permissions, requiredPermission) {
// 			return response.Error(c, fiber.StatusForbidden, "Permission denied: You do not have the required permission")
// 		}

// 		return c.Next()
// 	}
// }
