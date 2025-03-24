package main

import (
	"api-peak-form/internal/api"
	"api-peak-form/internal/config"
	"api-peak-form/internal/connection"
	"api-peak-form/internal/repository"
	"api-peak-form/internal/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cnf := config.Get()
	dbConnection := connection.GetDatabase(cnf.Database)

	app := fiber.New()

	//jwtMidd := jwtMid.New(jwtMid.Config{
	//	SigningKey: jwtMid.SigningKey{Key: []byte(cnf.Jwt.Key)},
	//	ErrorHandler: func(c *fiber.Ctx, err error) error {
	//		return c.Status(http.StatusUnauthorized).
	//			JSON(fiber.Map{
	//				"status":  "error",
	//				"message": "Invalid token",
	//			})
	//	},
	//})
	datadumy.AddDefaultUser(dbConnection)
	datadumy.AddExercise(dbConnection)
	datadumy.AddSchedules(dbConnection)
	datadumy.AddUserSchedules(dbConnection)

	uerRepository := repository.NewUserRepository(dbConnection)
	scheduleRepository := repository.NewSchedule(dbConnection)
	logRepository := repository.NewLogRepository(dbConnection)

	authService := service.NewAuthService(cnf, uerRepository)
	scheduleService := service.NewScheduleService(scheduleRepository)
	logService := service.NewLogService(logRepository)

	api.NewAuthApi(app, authService)
	api.NewScheduleApi(app, scheduleService)
	api.NewLogApi(app, logService)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
