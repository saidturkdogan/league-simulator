package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/user/league-simulator/src/config"
	"github.com/user/league-simulator/src/controller"
	"github.com/user/league-simulator/src/database"
	"github.com/user/league-simulator/src/docs"
	_ "github.com/user/league-simulator/src/docs" // Import for Swagger docs
	"github.com/user/league-simulator/src/middleware"
	"github.com/user/league-simulator/src/repository"
	"github.com/user/league-simulator/src/service"
)

// @title Football League Simulator API
// @version 1.0
// @description API for simulating a football league with 4 teams
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database")

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully")

	// Seed database with initial data
	if err := database.SeedData(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}
	log.Println("Database seeded successfully")

	// Create repositories
	repo := repository.NewPostgresRepository(db)

	// Create services
	svc := service.NewService(repo)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Football League Simulator",
		ErrorHandler: customErrorHandler,
	})

	// Setup middleware
	middleware.SetupMiddleware(app)

	// Setup Swagger
	docs.SetupSwagger(app)

	// Setup routes
	controller.SetupRoutes(app, svc)

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %d", cfg.Server.Port)
		if err := app.Listen(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Shutdown server
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

// customErrorHandler handles errors in a consistent way
func customErrorHandler(c *fiber.Ctx, err error) error {
	// Default status code is 500
	code := fiber.StatusInternalServerError

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Return JSON response
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
