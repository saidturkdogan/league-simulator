package controller

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/user/league-simulator/src/model"
	"github.com/user/league-simulator/src/service"
)

// TeamController handles HTTP requests for teams
type TeamController struct {
	service *service.TeamService
}

// NewTeamController creates a new TeamController
func NewTeamController(service *service.TeamService) *TeamController {
	return &TeamController{
		service: service,
	}
}

// GetTeams godoc
// @Summary Get all teams
// @Description Get a list of all teams in the league
// @Tags teams
// @Accept json
// @Produce json
// @Success 200 {array} model.Team
// @Failure 500 {object} ErrorResponse
// @Router /teams [get]
func (c *TeamController) GetTeams(ctx *fiber.Ctx) error {
	teams, err := c.service.GetAll(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(teams)
}

// GetTeam godoc
// @Summary Get a team by ID
// @Description Get a specific team by its ID
// @Tags teams
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Success 200 {object} model.Team
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /teams/{id} [get]
func (c *TeamController) GetTeam(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid team ID"})
	}

	team, err := c.service.GetByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(team)
}

// CreateTeam godoc
// @Summary Create a new team
// @Description Create a new team with the provided information
// @Tags teams
// @Accept json
// @Produce json
// @Param team body model.Team true "Team information"
// @Success 201 {object} model.Team
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /teams [post]
func (c *TeamController) CreateTeam(ctx *fiber.Ctx) error {
	var team model.Team
	if err := ctx.BodyParser(&team); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid request payload"})
	}

	if err := c.service.Create(ctx.Context(), &team); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(team)
}

// UpdateTeam godoc
// @Summary Update a team
// @Description Update a team with the provided information
// @Tags teams
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Param team body model.Team true "Team information"
// @Success 200 {object} model.Team
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /teams/{id} [put]
func (c *TeamController) UpdateTeam(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid team ID"})
	}

	var team model.Team
	if err := ctx.BodyParser(&team); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid request payload"})
	}

	team.ID = id
	if err := c.service.Update(ctx.Context(), &team); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(team)
}

// DeleteTeam godoc
// @Summary Delete a team
// @Description Delete a team by its ID
// @Tags teams
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /teams/{id} [delete]
func (c *TeamController) DeleteTeam(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid team ID"})
	}

	if err := c.service.Delete(ctx.Context(), id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(SuccessResponse{Result: "success"})
}

// CreateInitialTeams godoc
// @Summary Create initial teams
// @Description Create the initial 4 teams for the league
// @Tags teams
// @Accept json
// @Produce json
// @Success 201 {array} model.Team
// @Failure 500 {object} ErrorResponse
// @Router /teams/initialize [post]
func (c *TeamController) CreateInitialTeams(ctx *fiber.Ctx) error {
	log.Println("Creating initial teams...")
	
	// Check if teams already exist
	existingTeams, err := c.service.GetAll(ctx.Context())
	if err != nil {
		log.Printf("Error checking existing teams: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}
	
	if len(existingTeams) > 0 {
		log.Printf("Found %d existing teams, returning them", len(existingTeams))
		return ctx.Status(fiber.StatusOK).JSON(existingTeams)
	}
	
	// Create initial teams
	teams, err := c.service.CreateInitialTeams(ctx.Context())
	if err != nil {
		log.Printf("Error creating initial teams: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}
	
	log.Printf("Successfully created %d initial teams", len(teams))
	return ctx.Status(fiber.StatusCreated).JSON(teams)
}
