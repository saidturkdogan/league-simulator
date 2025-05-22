package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user/league-simulator/src/model"
	"github.com/user/league-simulator/src/service"
)

// MatchController handles HTTP requests for matches
type MatchController struct {
	service *service.MatchService
}

// NewMatchController creates a new MatchController
func NewMatchController(service *service.MatchService) *MatchController {
	return &MatchController{
		service: service,
	}
}

// GetMatches godoc
// @Summary Get all matches
// @Description Get a list of all matches or matches for a specific week
// @Tags matches
// @Accept json
// @Produce json
// @Param week query int false "Week number"
// @Success 200 {array} model.Match
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /matches [get]
func (c *MatchController) GetMatches(ctx *fiber.Ctx) error {
	// Check if week query parameter is provided
	weekStr := ctx.Query("week")
	if weekStr != "" {
		week := ctx.QueryInt("week", 0) // Default to 0 if conversion fails
		if week == 0 && weekStr != "0" {
			return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid week parameter"})
		}

		matches, err := c.service.GetByWeek(ctx.Context(), week)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
		}

		return ctx.JSON(matches)
	}

	// Get all matches
	matches, err := c.service.GetAll(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(matches)
}

// GetMatch godoc
// @Summary Get a match by ID
// @Description Get a specific match by its ID
// @Tags matches
// @Accept json
// @Produce json
// @Param id path int true "Match ID"
// @Success 200 {object} model.Match
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /matches/{id} [get]
func (c *MatchController) GetMatch(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid match ID"})
	}

	match, err := c.service.GetByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(match)
}

// CreateMatch godoc
// @Summary Create a new match
// @Description Create a new match with the provided information
// @Tags matches
// @Accept json
// @Produce json
// @Param match body model.Match true "Match information"
// @Success 201 {object} model.Match
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /matches [post]
func (c *MatchController) CreateMatch(ctx *fiber.Ctx) error {
	var match model.Match
	if err := ctx.BodyParser(&match); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid request payload"})
	}

	if err := c.service.Create(ctx.Context(), &match); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(match)
}

// UpdateMatch godoc
// @Summary Update a match
// @Description Update a match with the provided information
// @Tags matches
// @Accept json
// @Produce json
// @Param id path int true "Match ID"
// @Param match body model.Match true "Match information"
// @Success 200 {object} model.Match
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /matches/{id} [put]
func (c *MatchController) UpdateMatch(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid match ID"})
	}

	var match model.Match
	if err := ctx.BodyParser(&match); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid request payload"})
	}

	match.ID = id
	if err := c.service.Update(ctx.Context(), &match); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(match)
}
