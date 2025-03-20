package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func Get() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file: ", err.Error())
	}

	return &Config{
		Server: Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		Database: Database{
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			Name: os.Getenv("DB_NAME"),
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASSWORD"),
			Tz:   os.Getenv("DB_TZ"),
		},
	}
}
