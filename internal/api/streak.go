package api

import (
	"api-peak-form/domain"
	"github.com/gofiber/fiber/v2"
)

type streakApi struct {
	userService domain.UserService
}

func NewStreakApi(app *fiber.App, userService domain.UserService) {
	user := streakApi{
		userService: userService,
	}

	app.Get("/streak/check/:id", user.CheckStreak)
	app.Put("/streak/update/:id", user.UpdateStreak)
	app.Get("/leaderboard", user.GetLeaderboard)
}

func (ua streakApi) CheckStreak(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "User ID is required",
		})
	}

	streak, err := ua.userService.CheckStreak(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to check streak",
			"details": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"data":    streak,
		"message": "Streak checked successfully",
	})
}

func (ua streakApi) UpdateStreak(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "User ID is required",
		})
	}

	streak, err := ua.userService.UpdateStreak(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update streak",
			"details": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"data":    streak,
		"message": "Streak updated successfully",
	})
}

func (ua streakApi) GetLeaderboard(ctx *fiber.Ctx) error {
	users, err := ua.userService.GetAllUsersDesc(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve leaderboard",
			"details": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Leaderboard retrieved successfully",
		"data":    users,
	})
}
