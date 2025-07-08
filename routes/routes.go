package routes

import (
	"github.com/sabuj0338/go-task-manager/internal/auth"
	"github.com/sabuj0338/go-task-manager/internal/middleware"
	"github.com/sabuj0338/go-task-manager/internal/rbac"
	"github.com/sabuj0338/go-task-manager/internal/task"
	"github.com/sabuj0338/go-task-manager/internal/user"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(router fiber.Router) {
	router.Post("/register", auth.RegisterHandler)
	router.Post("/login", auth.LoginHandler)
	router.Post("/refresh-token", auth.RefreshTokenHandler)
	router.Post("/mfa/verify-code", auth.VerifyMFACodeHandler)

	// OTP-based password reset (for mobile)
	router.Post("/forgot-password", auth.ForgotPasswordHandler)
	router.Post("/reset-password", auth.ResetPasswordHandler)

	// Link-based password reset (for web)
	router.Post("/forgot-password-link", auth.ForgotPasswordLinkHandler)
	router.Post("/reset-password-link", auth.ResetPasswordWithTokenHandler)

	router.Post("/email/send", auth.SendEmailVerificationHandler)
	router.Post("/email/verify", auth.VerifyEmailCodeHandler)

	router.Post("/logout", auth.LogoutHandler)

	protected := router.Group("/", middleware.JWTProtected())
	protected.Post("/mfa/setup", auth.SetupMFAHandler)
	protected.Post("/mfa/verify", auth.VerifyMFAHandler)
	protected.Post("/mfa/disable", auth.DisableMFAHandler)
}

func RegisterUserRoutes(router fiber.Router) {
	router.Use(middleware.JWTProtected())

	// Routes requiring specific permissions for user management
	router.Get("/", middleware.RequirePermission("users:read"), user.GetUsers)
	router.Post("/", middleware.RequirePermission("users:create"), user.CreateUser)
	router.Delete("/:id", middleware.RequirePermission("users:delete"), user.DeleteUser)

	// A user can view/update their own profile.
	// An admin with 'users:read' or 'users:update' can also access these.
	// The logic for this check should ideally be inside the handlers, where you
	// can compare the authenticated user's ID with the ID from the URL parameter.
	router.Get("/:id", user.GetUser)
	router.Put("/:id", user.UpdateUser)
}

func RegisterTaskRoutes(router fiber.Router) {
	// All task routes require authentication first.
	// Then, each route is protected by a specific permission.
	router.Use(middleware.JWTProtected())

	router.Post("/", middleware.RequirePermission("tasks:create"), task.CreateTaskHandler)
	router.Get("/", middleware.RequirePermission("tasks:read"), task.GetTasksHandler)
	router.Get("/:id", middleware.RequirePermission("tasks:read"), task.GetTaskHandler)
	router.Put("/:id", middleware.RequirePermission("tasks:update"), task.UpdateTaskHandler)
	router.Delete("/:id", middleware.RequirePermission("tasks:delete"), task.DeleteTaskHandler)
}

func RegisterRBACRoutes(router fiber.Router) {
	// All RBAC routes are protected and require specific high-level permissions.
	// For simplicity, we can reuse 'users:update' as a proxy for these admin actions.
	// In a larger system, you might create specific permissions like 'roles:create', 'permissions:assign', etc.
	router.Use(middleware.JWTProtected(), middleware.RequirePermission("users:update"))

	router.Post("/roles", rbac.CreateRoleHandler)
	router.Post("/permissions", rbac.CreatePermissionHandler)
	router.Post("/roles/assign-permission", rbac.AssignPermissionToRoleHandler)
	router.Post("/users/assign-role", rbac.AssignRoleToUserHandler)
}
