package api

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"github.com/gofiber/fiber/v2"
)

type logApi struct {
	logService domain.LogService
}

func NewLogApi(app *fiber.App, logService domain.LogService) {
	la := logApi{logService: logService}

	app.Post("/logs", la.Create)
}

func (la logApi) Create(ctx *fiber.Ctx) error {
	var req dto.LogRequest

	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body. Ensure all fields are correctly formatted",
		})
	}

	if req.UserID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "UserID is required",
		})
	}

	if req.ExerciseID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ExerciseID is required",
		})
	}

	if req.Set <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Set must be a positive integer",
		})
	}

	if req.Repetition <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Repetition must be a positive integer",
		})
	}

	err = la.logService.Create(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create log. Please try again later",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Log created successfully",
	})
}
