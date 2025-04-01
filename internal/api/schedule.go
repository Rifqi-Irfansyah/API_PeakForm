package api

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"api-peak-form/internal/util"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type scheduleApi struct {
	scheduleService domain.ScheduleService
}

func NewScheduleApi(app *fiber.App, scheduleService domain.ScheduleService) {
	sa := scheduleApi{
		scheduleService: scheduleService,
	}

	app.Get("/schedule", sa.FindByUID)
	app.Post("/schedule", sa.Create)
	app.Put("/schedule", sa.Update)
	app.Put("/schedule/exercise", sa.UpdateExerciseList)
	app.Delete("/schedule", sa.Delete)
	app.Delete("/schedule/exercise", sa.DeleteExercise)
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

func (sa scheduleApi) Update(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.UpdateScheduleRequest
	req.ID = ctx.Query("id")
	req.Day = ctx.QueryInt("day")

	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation failed",
			"details": fails,
		})
	}

	err := sa.scheduleService.Update(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    req,
	})
}

func (sa scheduleApi) UpdateExerciseList(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.UpdateExerciseScheduleRequest
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

	err := sa.scheduleService.UpdateExerciseSchedule(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    req,
	})
}

func (sa scheduleApi) FindByUID(ctx *fiber.Ctx) error {
	_, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	UID := ctx.Query("UID")
	if UID == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "UID parameter is required",
		})
	}

	res, err := sa.scheduleService.FindByUID(ctx.Context(), UID)
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

func (sa scheduleApi) Delete(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id := ctx.Query("id_schedule")

	idUintSchedule, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		println("Error saat membaca ID API Req: ", err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	err = sa.scheduleService.DeleteSchedule(c, uint(idUintSchedule))
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"error": "schedule not found"})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "schedule deleted successfully",
	})
}

func (sa scheduleApi) DeleteExercise(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id 			:= ctx.Query("id_schedule")
	id_exercise	:= ctx.Query("id_exercise")

	uintId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		println("Error saat membaca ID API Req: ", err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	intExerciseID, err := strconv.Atoi(id_exercise)
	if err != nil {
		println("Error saat membaca ID Exercise API Req: ", err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Exercise ID format",
		})
	}

	err = sa.scheduleService.DeleteExerciseSchedule(c, uint(uintId), int(intExerciseID))
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Exercise Schedule not found"})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "schedule deleted successfully",
	})
}