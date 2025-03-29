package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func Get() *Config {
	err := godotenv.Load()

	if err != nil {
		logrus.Fatal("Error loading .env file: ", err.Error())
	}

	expInt, err := strconv.Atoi(os.Getenv("JWT_EXP"))
	if err != nil {
		logrus.Fatal("Invalid JWT_EXP value: ", err.Error())
	}

	config := &Config{
		Server: Server{
			Host: getEnv("SERVER_HOST"),
			Port: getEnv("SERVER_PORT"),
		},
		Database: Database{
			Host: getEnv("DB_HOST"),
			Port: getEnv("DB_PORT"),
			Name: getEnv("DB_NAME"),
			User: getEnv("DB_USER"),
			Pass: getEnv("DB_PASSWORD"),
			Tz:   getEnv("DB_TZ"),
		},
		Jwt: Jwt{
			Key: getEnv("JWT_KEY"),
			Exp: expInt,
		},
	}

	logrus.Infof("Configuration loaded: %+v", config)
	return config
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logrus.Fatalf("Environment variable %s is not set", key)
	}
	return value
}
