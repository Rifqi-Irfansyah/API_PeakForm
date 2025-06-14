package api

import (
    "api-peak-form/domain"
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
}

func (u userApi) UploadPhoto(c *fiber.Ctx) error {
    id := c.Params("id")
	baseURL := os.Getenv("BASE_URL")
	fmt.Println("DEBUG BASE_URL:", baseURL)

	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}

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

    ext := strings.ToLower(filepath.Ext(file.Filename)) // Normalisasi ke huruf kecil
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
    savePath := fmt.Sprintf("public/profile/%s", filename)

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

    url := fmt.Sprintf("%s/profile/%s", baseURL, filename)

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
	baseURL := os.Getenv("BASE_URL")

	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	
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
			photoURL = fmt.Sprintf("%s/profile/%s%s", baseURL, id, ext)
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