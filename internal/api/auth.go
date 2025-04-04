package api

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"api-peak-form/internal/util"
	"context"
	"errors"
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
	app.Post("/register", aa.Register)
	app.Post("/register/verify-otp", aa.VerifyRegisterOTP)
	app.Post("/auth/forgot-password", aa.ForgotPassword)
	app.Post("/auth/reset-password", aa.ResetPassword)
}

func (aa authApi) Login(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.AuthRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validation failed",
			"details": fails,
		})
	}

	res, err := aa.authService.Login(c, req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return ctx.Status(http.StatusRequestTimeout).JSON(fiber.Map{
				"status":  "error",
				"message": "Request timed out",
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Login failed",
			"details": err.Error(),
		})
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Login successful",
		"data":    res,
	})
}

func (aa authApi) Register(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validation failed",
			"details": fails,
		})
	}

	err := aa.authService.Register(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Registration failed",
			"details": err.Error(),
		})
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Registration successful, OTP sent",
	})
}

func (aa authApi) VerifyRegisterOTP(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.VerifyOTPRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validation failed",
			"details": fails,
		})
	}

	err := aa.authService.VerifyRegisterOTP(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "OTP verification failed",
			"details": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "OTP verified, registration complete",
	})
}

func (aa authApi) ForgotPassword(ctx *fiber.Ctx) error {
	var req dto.ForgotPasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validation failed",
			"details": fails,
		})
	}

	err := aa.authService.ForgotPassword(ctx.Context(), req.Email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to send OTP",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "OTP sent successfully",
	})
}

func (aa authApi) ResetPassword(ctx *fiber.Ctx) error {
	var req dto.ResetPasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validation failed",
			"details": fails,
		})
	}

	err := aa.authService.ResetPassword(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Password updated successfully",
	})
}
