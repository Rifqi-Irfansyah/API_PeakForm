package api

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"api-peak-form/internal/util"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type userApi struct {
	userService domain.UserService
}

func NewUserApi(app *fiber.App, userService domain.UserService) {
	user := userApi{
		userService: userService,
	}

	app.Put("/user/streak", user.CheckStreak)
	app.Put("/user/point", user.UpdateStreak)
}

func (ua userApi) CheckStreak(ctx *fiber.Ctx) error {
	var req dto.UserStreakRequest

	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"details": err.Error(),
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
			"status":  "error",
			"message": "Validation failed",
			"details": validationErrors,
		})
	}

	streak, err := ua.userService.CheckStreak(ctx.Context(), req.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to check streak",
			"details": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"streak":  streak,
		"message": "Streak checked successfully",
	})
}

func (ua userApi) UpdateStreak(ctx *fiber.Ctx) error {
	var req dto.UserStreakRequest

	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"details": err.Error(),
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
			"status":  "error",
			"message": "Validation failed",
			"details": validationErrors,
		})
	}

	streak, err := ua.userService.UpdateStreak(ctx.Context(), req.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update streak",
			"details": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"streak":  streak,
		"message": "Streak updated successfully",
	})
}
