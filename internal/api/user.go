package api

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"github.com/gofiber/fiber/v2"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type userApi struct {
	userService domain.UserService
}

func NewUserApi(app *fiber.App, userService domain.UserService) {
	user := userApi{
		userService: userService,
	}

	app.Post("/users/:id/photo", user.UploadPhoto)
	app.Get("/users/:id/photo", user.GetPhoto)
	app.Get("/users/:id", user.FindByID)
	app.Static("/users/static/photo", "./public/profile")
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

    ext := strings.ToLower(filepath.Ext(file.Filename))
    if ext == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "File has no extension",
        })
    }

    allowedExt := []string{".jpg", ".jpeg", ".png"}
    if !contains(allowedExt, ext) {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": fmt.Sprintf("Unsupported file extension. Allowed: %v", allowedExt),
        })
    }

    filename := fmt.Sprintf("%s%s", id, ext)
    savePath := filepath.Join("public", "profile", filename)

    for _, e := range allowedExt {
        oldPath := filepath.Join("public", "profile", fmt.Sprintf("%s%s", id, e))
        if oldPath != savePath {
            _ = os.Remove(oldPath)
        }
    }

    if err := c.SaveFile(file, savePath); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Failed to save file",
            "details": err.Error(),
        })
    }

    err = u.userService.UpdatePhoto(c.Context(), id, filename)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Failed to update photo",
            "details": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status":  "success",
        "message": "Photo uploaded successfully",
        "data": fiber.Map{
            "filename": filename,
        },
    })
}


func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

func (u userApi) GetPhoto(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "User ID is required",
		})
	}

	filename, err := u.userService.GetPhotoFilename(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Photo not found",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Photo fetched successfully",
		"data": fiber.Map{
			"filename": filename,
		},
	})
}

func (u userApi) FindByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "User ID is required",
		})
	}

	user, err := u.userService.FindByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find user",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User found successfully",
		"data": dto.UserResponse{
			Email:    user.Email,
			Name:     user.Name,
			Point:    user.Point,
			Streak:   user.Streak,
			PhotoURL: user.PhotoURL,
		},
	})
}
