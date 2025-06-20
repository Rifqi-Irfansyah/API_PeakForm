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
	logService      domain.LogService
	userService     domain.UserService
	exerciseService domain.ExerciseService
}

func NewLogApi(app *fiber.App, logService domain.LogService, userService domain.UserService, exerciseService domain.ExerciseService) {
	la := logApi{logService: logService, userService: userService, exerciseService: exerciseService}

	app.Post("/logs/create", la.Create)
	app.Get("/logs/:id", la.FindByUserID)
	app.Get("/logs/:id/summary", la.GetUserWorkoutSummary)
	app.Get("/logs/is-exercised-today/:id", la.HasUserExercisedToday)

}

func (la logApi) Create(ctx *fiber.Ctx) error {
	var req dto.LogRequest

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

	exercise, err := la.exerciseService.GetExerciseByID(ctx.Context(), req.ExerciseID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch exercise details",
			"details": err.Error(),
		})
	}

	points, err := la.userService.UpdatePoint(ctx.Context(), req.UserID, exercise.Difficulty, req.Repetition, req.Set)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update user points",
			"details": err.Error(),
		})
	}

	err = la.logService.Create(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create log",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Log created successfully",
		"points":  points,
		"data":    req,
	})
}

func (la logApi) FindByUserID(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")
	if userID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "UserID is required",
		})
	}

	logs, err := la.logService.FindByUserID(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch logs for user ID " + userID,
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Logs fetched successfully",
		"data":    logs,
	})
}

func (la logApi) GetUserWorkoutSummary(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")
	if userID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "UserID is required",
		})
	}

	summary, err := la.logService.GetUserWorkoutSummary(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch workout summary for user ID " + userID,
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Workout summary fetched successfully",
		"data":    summary,
	})
}

func (la logApi) HasUserExercisedToday(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")
	if userID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "UserID is required",
		})
	}

	hasExercised, err := la.logService.HasUserExercisedToday(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to check if user has exercised today",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User exercise status fetched successfully",
		"data":    hasExercised,
	})
}
