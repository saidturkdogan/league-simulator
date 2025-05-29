package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user/league-simulator/src/service"
)

// PredictionController handles HTTP requests for predictions
type PredictionController struct {
	service *service.PredictionService
}

// NewPredictionController creates a new PredictionController
func NewPredictionController(service *service.PredictionService) *PredictionController {
	return &PredictionController{
		service: service,
	}
}

// PredictFinalStandings godoc
// @Summary Predict final standings
// @Description Predict the final standings for a league after all weeks are played
// @Tags predictions
// @Accept json
// @Produce json
// @Param id path int true "League ID"
// @Success 200 {object} model.Standings
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leagues/{id}/predict [get]
func (c *PredictionController) PredictFinalStandings(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid league ID"})
	}

	standings, err := c.service.PredictFinalStandings(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(standings)
}

// GetPredictionWithConfidence godoc
// @Summary Get detailed predictions with confidence levels
// @Description Get comprehensive predictions with probabilities and confidence levels (available after week 4)
// @Tags predictions
// @Accept json
// @Produce json
// @Param id path int true "League ID"
// @Success 200 {object} model.PredictionResult
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leagues/{id}/predictions [get]
func (c *PredictionController) GetPredictionWithConfidence(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Geçersiz liga ID"})
	}

	predictions, err := c.service.GetPredictionWithConfidence(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(predictions)
}
