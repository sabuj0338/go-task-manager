package routes

import (
	"github.com/sabuj0338/go-task-manager/internal/auth"
	"github.com/sabuj0338/go-task-manager/internal/middleware"
	"github.com/sabuj0338/go-task-manager/internal/task"
	"github.com/sabuj0338/go-task-manager/internal/user"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(router fiber.Router) {
	router.Post("/register", auth.RegisterHandler)
	router.Post("/login", auth.LoginHandler)
	router.Post("/refresh-token", auth.RefreshTokenHandler)
	router.Post("/mfa/verify-code", auth.VerifyMFACodeHandler)

	router.Post("/forgot-password", auth.ForgotPasswordHandler)
	router.Post("/reset-password", auth.ResetPasswordHandler)

	router.Post("/email/send", auth.SendEmailVerificationHandler)
	router.Post("/email/verify", auth.VerifyEmailCodeHandler)

	protected := router.Group("/", middleware.JWTProtected())
	protected.Post("/mfa/setup", auth.SetupMFAHandler)
	protected.Post("/mfa/verify", auth.VerifyMFAHandler)
}

func RegisterUserRoutes(router fiber.Router) {
	router.Use(middleware.JWTProtected())

	// router.Get("/", user.GetUsers)
	// router.Get("/:id", user.GetUser)
	// router.Put("/:id", user.UpdateUser)
	// router.Delete("/:id", user.DeleteUser)

	// Only admin can list or delete users
	router.Get("/", middleware.RequireRole("admin"), user.GetUsers)
	router.Delete("/:id", middleware.RequireRole("admin"), user.DeleteUser)

	// Admin or user can view/update themselves
	router.Get("/:id", user.GetUser)
	router.Put("/:id", user.UpdateUser)
}

func RegisterTaskRoutes(router fiber.Router) {
	router.Use(middleware.JWTProtected())

	router.Post("/", task.CreateTaskHandler)
	router.Get("/", task.GetTasksHandler)
	router.Get("/:id", task.GetTaskHandler)
	router.Put("/:id", task.UpdateTaskHandler)
	router.Delete("/:id", task.DeleteTaskHandler)
}

func SetupRoutes(app *fiber.App) {
	api := app.Group("/v1/api")

	// Auth routes
	auth := api.Group("/auth")
	RegisterAuthRoutes(auth)

	// User routes
	users := api.Group("/users")
	RegisterUserRoutes(users)

	// Task routes
	tasks := api.Group("/tasks")
	RegisterTaskRoutes(tasks)
}
