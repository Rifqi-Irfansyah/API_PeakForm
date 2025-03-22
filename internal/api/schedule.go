package api

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"context"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

type scheduleApi struct {
	scheduleService domain.ScheduleService
}

func NewScheduleApi(app *fiber.App, scheduleService domain.ScheduleService) {
	aa := scheduleApi{
		scheduleService: scheduleService,
	}

	app.Get("/schedule", aa.FindByUID)
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
