package main

import (
	"log"
	"os"
	"time"

	"github.com/sabuj0338/go-task-manager/pkg/database"
	"github.com/sabuj0338/go-task-manager/pkg/redis"
	"github.com/sabuj0338/go-task-manager/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found, relying on system environment")
	}

	app := fiber.New()

	// Init Redis + MySQL
	redis.ConnectRedis()
	database.ConnectMySQL()

	// -------------------
	// ✅ Middlewares
	// -------------------

	app.Use(recover.New()) // recover from panics
	app.Use(logger.New())  // log all requests
	app.Use(helmet.New())  // secure headers

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // You can restrict by domain (e.g., https://example.com)
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        100,             // 100 requests
		Expiration: 1 * time.Minute, // per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests. Please try again later.",
			})
		},
	}))

	// Routes
	api := app.Group("/v1/api")
	routes.RegisterAuthRoutes(api.Group("/auth"))
	routes.RegisterUserRoutes(api.Group("/users"))
	routes.RegisterTaskRoutes(api.Group("/tasks"))

	// Static files (optional)
	app.Static("/public", "./public")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
