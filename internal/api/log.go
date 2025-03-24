package api

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"api-peak-form/internal/util"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type logApi struct {
	logService domain.LogService
}

func NewLogApi(app *fiber.App, logService domain.LogService) {
	la := logApi{logService: logService}

	app.Post("/logs", la.Create)
	app.Get("/logs/:id", la.FindByUserID)
}

func (la logApi) Create(ctx *fiber.Ctx) error {
	var req dto.LogRequest

	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body. Ensure all fields are correctly formatted",
		})
	}

	err = validator.New().Struct(req)
	if err != nil {
		validationErrors := make(map[string]string)

		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			for _, e := range errs {
				validationErrors[e.StructField()] = util.TranslateTag(e)
			}
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body. Ensure all fields are correctly formatted",
			"details": validationErrors,
		})
	}

	err = la.logService.Create(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create log",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Log created successfully",
		"data":    req,
	})
}

func (la logApi) FindByUserID(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")
	if userID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "UserID is required",
		})
	}

	logs, err := la.logService.FindByUserID(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch logs. Please try again later",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": logs,
	})
}
