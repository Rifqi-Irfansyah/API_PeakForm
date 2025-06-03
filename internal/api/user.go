package api

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"api-peak-form/internal/util"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"fmt"
	"os"
	"path/filepath"
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
	app.Post("/user/:id/photo", user.UploadPhoto)
	app.Get("/user/:id/photo", user.GetPhoto)
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

func (u userApi) UploadPhoto(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "User ID is required",
		})
	}

	file, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "File not found",
			"details": err.Error(),
		})
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", id, ext)
	savePath := fmt.Sprintf("public/profile/%s", filename)

	allowedExt := []string{".jpg", ".jpeg", ".png"}
	for _, ext := range allowedExt {
		oldPath := fmt.Sprintf("public/profile/%s%s", id, ext)
		if oldPath != savePath {
			_ = os.Remove(oldPath) 
		}
	}

	err = c.SaveFile(file, savePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to save file",
			"details": err.Error(),
		})
	}

	url := fmt.Sprintf("http://localhost:3000/profile/%s", filename)

	err = u.userService.UpdatePhoto(c.Context(), id, url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update photo URL",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Photo uploaded successfully",
		"data": fiber.Map{
			"url": url,
		},
	})
}

func (u userApi) GetPhoto(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "User ID is required",
		})
	}

	allowedExt := []string{".jpg", ".jpeg", ".png"}
	var photoURL string
	found := false

	for _, ext := range allowedExt {
		filename := fmt.Sprintf("public/profile/%s%s", id, ext)
		if _, err := os.Stat(filename); err == nil {
			photoURL = fmt.Sprintf("http://localhost:3000/profile/%s%s", id, ext)
			found = true
			break
		}
	}

	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Photo not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Photo fetched successfully",
		"data": fiber.Map{
			"url": photoURL,
		},
	})
}
