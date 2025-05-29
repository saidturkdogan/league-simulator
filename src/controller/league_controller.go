package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user/league-simulator/src/service"
)

// LeagueController handles HTTP requests for leagues
type LeagueController struct {
	service *service.LeagueService
}

// NewLeagueController creates a new LeagueController
func NewLeagueController(service *service.LeagueService) *LeagueController {
	return &LeagueController{
		service: service,
	}
}

// CreateLeague godoc
// @Summary Create a new league
// @Description Create a new league with the provided name
// @Tags leagues
// @Accept json
// @Produce json
// @Param league body CreateLeagueRequest true "League information"
// @Success 201 {object} model.League
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leagues [post]
func (c *LeagueController) CreateLeague(ctx *fiber.Ctx) error {
	var request CreateLeagueRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid request payload"})
	}

	if request.Name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "League name is required"})
	}

	league, err := c.service.Create(ctx.Context(), request.Name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(league)
}

// GetLeague godoc
// @Summary Get a league by ID
// @Description Get a specific league by its ID
// @Tags leagues
// @Accept json
// @Produce json
// @Param id path int true "League ID"
// @Success 200 {object} model.League
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /leagues/{id} [get]
func (c *LeagueController) GetLeague(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid league ID"})
	}

	league, err := c.service.GetByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(league)
}

// SimulateWeek godoc
// @Summary Simulate a week of matches
// @Description Simulate all matches for the next week in the league
// @Tags leagues
// @Accept json
// @Produce json
// @Param id path int true "League ID"
// @Success 200 {object} model.Standings
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leagues/{id}/simulate [post]
func (c *LeagueController) SimulateWeek(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid league ID"})
	}

	standings, err := c.service.SimulateWeek(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(standings)
}

// GetStandings godoc
// @Summary Get current standings
// @Description Get the current standings for a league
// @Tags leagues
// @Accept json
// @Produce json
// @Param id path int true "League ID"
// @Success 200 {object} model.Standings
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leagues/{id}/standings [get]
func (c *LeagueController) GetStandings(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid league ID"})
	}

	standings, err := c.service.GetCurrentStandings(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(standings)
}

// SimulateAllWeeks - Tüm kalan haftaları simüle et
// @Summary Tüm kalan haftaları simüle et
// @Description Liga bitene kadar otomatik olarak tüm haftaları simüle eder
// @Tags leagues
// @Accept json
// @Produce json
// @Param id path int true "Liga ID"
// @Success 200 {object} model.LeagueSimulationResult
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leagues/{id}/simulate-all [post]
func (c *LeagueController) SimulateAllWeeks(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Geçersiz liga ID"})
	}

	result, err := c.service.SimulateAllRemainingWeeks(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(result)
}

// GetWeeklyMatches - Haftalık maçları getir
// @Summary Belirli bir haftanın maçlarını getir
// @Description Ligada belirli bir haftanın tüm maçlarını getir
// @Tags leagues
// @Accept json
// @Produce json
// @Param id path int true "Liga ID"
// @Param week path int true "Hafta numarası"
// @Success 200 {array} model.Match
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leagues/{id}/weeks/{week}/matches [get]
func (c *LeagueController) GetWeeklyMatches(ctx *fiber.Ctx) error {
	leagueID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Geçersiz liga ID"})
	}

	week, err := ctx.ParamsInt("week")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Geçersiz hafta numarası"})
	}

	matches, err := c.service.GetWeeklyMatches(ctx.Context(), leagueID, week)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(matches)
}
