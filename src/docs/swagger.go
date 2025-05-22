package docs

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
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

// Note: Swagger info is now defined in the docs package

// SetupSwagger sets up the Swagger documentation
func SetupSwagger(app *fiber.App) {
	// Register Swagger docs
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/swagger/doc.json",
		DeepLinking: true,
	}))

	// Serve the Swagger JSON file
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./src/docs/swagger.json")
	})
}
