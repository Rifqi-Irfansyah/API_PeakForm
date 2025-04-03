package api

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"api-peak-form/internal/util"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ExerciseAPI struct {
	service domain.ExerciseService
}

func NewExerciseAPI(app *fiber.App, service domain.ExerciseService) *ExerciseAPI {
	api := &ExerciseAPI{service}

	app.Post("/exercises", api.CreateExercise)
	app.Get("/exercises", api.GetExercises)
	app.Get("/exercises/:id", api.GetExerciseByID)
	app.Put("/exercises/:id", api.UpdateExercise)
	app.Delete("/exercises/:id", api.DeleteExercise)
	app.Static("/static/exercises", "./assets/exercises")

	return api
}

const baseURL = "http://localhost:3000"

func (api *ExerciseAPI) CreateExercise(c *fiber.Ctx) error {
	ctx := context.Background()

	var req dto.CreateExerciseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validation failed",
			"details": fails,
		})
	}

	file, err := c.FormFile("gif")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "GIF file is required",
		})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".gif" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid file format. Only GIF is allowed",
		})
	}

	savePath := filepath.Join("assets", "exercises", file.Filename)
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to save GIF file",
		})
	}

	req.Image = fmt.Sprintf("%s/static/exercises/%s", baseURL, file.Filename)

	if err := api.service.CreateExercise(ctx, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create exercise",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Exercise created successfully",
		"data": fiber.Map{
			"name":         req.Name,
			"type":         req.Type,
			"muscle":       req.Muscle,
			"equipment":    req.Equipment,
			"difficulty":   req.Difficulty,
			"instructions": req.Instructions,
			"image":          req.Image,
		},
	})
}

func (api *ExerciseAPI) GetExercises(c *fiber.Ctx) error {
	ctx := context.Background()

	exercises, err := api.service.GetExercises(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get exercises",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Exercises retrieved successfully",
		"data":    exercises,
	})
}

func (api *ExerciseAPI) GetExerciseByID(c *fiber.Ctx) error {
	ctx := context.Background()

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid exercise ID",
		})
	}

	exercise, err := api.service.GetExerciseByID(ctx, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Exercise not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Exercise retrieved successfully",
		"data":    exercise,
	})
}

func (api *ExerciseAPI) UpdateExercise(c *fiber.Ctx) error {
	ctx := context.Background()

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid exercise ID",
		})
	}

	var req dto.UpdateExerciseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validation failed",
			"details": fails,
		})
	}

	req.ID = uint(id)

	if err := api.service.UpdateExercise(ctx, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update exercise",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Exercise updated successfully",
	})
}

func (api *ExerciseAPI) DeleteExercise(c *fiber.Ctx) error {
	ctx := context.Background()

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid exercise ID",
		})
	}

	if err := api.service.DeleteExercise(ctx, uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete exercise",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Exercise deleted successfully",
	})
}
