package main

import (
	"api-peak-form/domain"
	"api-peak-form/internal/api"
	"api-peak-form/internal/config"
	"api-peak-form/internal/connection"
	"api-peak-form/internal/connection/datadumy"
	"api-peak-form/internal/repository"
	"api-peak-form/internal/service"
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	// jwtMiddleware "github.com/gofiber/jwt/v3"
	"net/http"
)

func main() {
	cnf := config.Get()
	dbConnection := connection.GetDatabase(cnf.Database)

	app := fiber.New()
	app.Static("/", "./assets") 

	err := dbConnection.AutoMigrate(&domain.User{})
	if err != nil {
		logrus.Fatal("Failed to migrate database:", err)
	}
	logrus.Info("Database migrated successfully")

	//jwtMid := jwtMiddleware.New(jwtMiddleware.Config{
	//	SigningKey:   []byte(cnf.Jwt.Key),
	//	ErrorHandler: jwtError,
	//})
	datadumy.AddDefaultUser(dbConnection)
	datadumy.AddExercise(dbConnection)
	datadumy.AddSchedules(dbConnection)
	//datadumy.AddUserSchedules(dbConnection)

	otpRepository := repository.NewOTPRepository()
	userRepository := repository.NewUserRepository(dbConnection)
	scheduleRepository := repository.NewSchedule(dbConnection)
	exerciseRepo := repository.NewExercise(dbConnection)
	logRepository := repository.NewLogRepository(dbConnection)

	scheduleService := service.NewScheduleService(scheduleRepository)
	exerciseService := service.NewExerciseService(exerciseRepo)
	authService := service.NewAuthService(cnf, userRepository, otpRepository)
	logService := service.NewLogService(logRepository, userRepository)
	userService := service.NewUserService(userRepository)

	// endpoints that do not require a token
	api.NewAuthApi(app, authService)

	// endpoints that require a token
	//app.Use(jwtMid)
	api.NewScheduleApi(app, scheduleService)
	api.NewLogApi(app, logService, userService, exerciseService) 
	api.NewExerciseAPI(app, exerciseService)

	logrus.Infof("Starting server at %s:%s", cnf.Server.Host, cnf.Server.Port)
	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}

func jwtError(c *fiber.Ctx, _ error) error {
	logrus.Warn("Invalid or expired token")
	return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
		"status":  "error",
		"message": "Invalid or expired token",
	})
}
