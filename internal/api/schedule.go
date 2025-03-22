package api

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"api-peak-form/internal/util"
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type scheduleApi struct {
	scheduleService domain.ScheduleService
}

func NewScheduleApi(app *fiber.App, scheduleService domain.ScheduleService) {
	aa := scheduleApi{
		scheduleService: scheduleService,
	}

	app.Get("/schedule", aa.FindByUID)
	app.Post("/schedule", aa.Create)
}

func (sa scheduleApi) Create(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.CreateScheduleRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation failed",
			"details": fails,
		})
	}

	err := sa.scheduleService.Create(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    req,
	})
}

func (aa scheduleApi) FindByUID(ctx *fiber.Ctx) error {
	_, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.ScheduleRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	res, err := aa.scheduleService.FindByUID(ctx.Context(), req.UID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Schedule Founded",
		"data":    res,
	})
}
