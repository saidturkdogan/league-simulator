package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/user/league-simulator/src/service"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Result string `json:"result"`
}

// CreateLeagueRequest represents a request to create a league
type CreateLeagueRequest struct {
	Name string `json:"name"`
}

// SetupRoutes sets up all the routes for the application
func SetupRoutes(app *fiber.App, service *service.Service) {
	// Create controllers
	teamController := NewTeamController(service.Team)
	matchController := NewMatchController(service.Match)
	leagueController := NewLeagueController(service.League)
	predictionController := NewPredictionController(service.Prediction)

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Swagger documentation is now set up in the docs package

	// API routes
	api := app.Group("/api")

	// Team routes
	teams := api.Group("/teams")
	teams.Get("/", teamController.GetTeams)
	teams.Get("/:id", teamController.GetTeam)
	teams.Post("/", teamController.CreateTeam)
	teams.Put("/:id", teamController.UpdateTeam)
	teams.Delete("/:id", teamController.DeleteTeam)
	teams.Post("/initialize", teamController.CreateInitialTeams)

	// Match routes
	matches := api.Group("/matches")
	matches.Get("/", matchController.GetMatches)
	matches.Get("/:id", matchController.GetMatch)
	matches.Post("/", matchController.CreateMatch)
	matches.Put("/:id", matchController.UpdateMatch)

	// League routes
	leagues := api.Group("/leagues")
	leagues.Post("/", leagueController.CreateLeague)
	leagues.Get("/:id", leagueController.GetLeague)
	leagues.Post("/:id/simulate", leagueController.SimulateWeek)
	leagues.Post("/:id/simulate-all", leagueController.SimulateAllWeeks)
	leagues.Get("/:id/standings", leagueController.GetStandings)
	leagues.Get("/:id/weeks/:week/matches", leagueController.GetWeeklyMatches)

	// Prediction routes
	leagues.Get("/:id/predict", predictionController.PredictFinalStandings)
	leagues.Get("/:id/predictions", predictionController.GetPredictionWithConfidence)

	// For backward compatibility, also add routes without /api prefix
	// Team routes
	app.Get("/teams", teamController.GetTeams)
	app.Get("/teams/:id", teamController.GetTeam)
	app.Post("/teams", teamController.CreateTeam)
	app.Put("/teams/:id", teamController.UpdateTeam)
	app.Delete("/teams/:id", teamController.DeleteTeam)
	app.Post("/teams/initialize", teamController.CreateInitialTeams)

	// Match routes
	app.Get("/matches", matchController.GetMatches)
	app.Get("/matches/:id", matchController.GetMatch)
	app.Post("/matches", matchController.CreateMatch)
	app.Put("/matches/:id", matchController.UpdateMatch)

	// League routes
	app.Post("/leagues", leagueController.CreateLeague)
	app.Get("/leagues/:id", leagueController.GetLeague)
	app.Post("/leagues/:id/simulate", leagueController.SimulateWeek)
	app.Post("/leagues/:id/simulate-all", leagueController.SimulateAllWeeks)
	app.Get("/leagues/:id/standings", leagueController.GetStandings)
	app.Get("/leagues/:id/weeks/:week/matches", leagueController.GetWeeklyMatches)

	// Prediction routes
	app.Get("/leagues/:id/predict", predictionController.PredictFinalStandings)
	app.Get("/leagues/:id/predictions", predictionController.GetPredictionWithConfidence)
}
