package api

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"context"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

type authApi struct {
	authService domain.AuthService
}

func NewAuthApi(app *fiber.App, authService domain.AuthService) {
	aa := authApi{
		authService: authService,
	}

	app.Post("/auth", aa.Login)
}

func (aa authApi) Login(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.AuthRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	res, err := aa.authService.Login(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Login success",
		"data":    res,
	})
}
