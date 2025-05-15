package api

import (
	"api-peak-form/domain"
	"github.com/gofiber/fiber/v2"
)

type statsApi struct {
	statsService domain.StatsService
}

func NewStatsApi(app *fiber.App, statsService domain.StatsService) {
	api := statsApi{statsService: statsService}
	app.Get("/stats/:id", api.GetUserStats)
}

func (s statsApi) GetUserStats(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")
	if userID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "UserID is required",
		})
	}

	summary, err := s.statsService.GetStatsByUserID(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch stats for user ID",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Stats fetched successfully",
		"data": summary,
	})
}