package migration

import (
	"api-peak-form/domain"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strings"
)

func Migrate(dbConnection *gorm.DB) error {
	if err := createEnums(dbConnection); err != nil {
		log.Printf("Warning: Failed to create enums: %v", err)
	}

	err := dbConnection.AutoMigrate(&domain.User{}, &domain.Log{}, &domain.Exercise{}, &domain.Schedule{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrated successfully")
	return nil
}

func createEnums(db *gorm.DB) error {
	enumQueries := []struct {
		Type  string
		Query string
	}{
		{"exercise_type", "CREATE TYPE exercise_type AS ENUM ('strength', 'cardio');"},
		{"muscle_group", "CREATE TYPE muscle_group AS ENUM ('abdominals', 'biceps', 'calves', 'chest', 'forearms', 'lats', 'lower_back', 'middle_back', 'neck', 'quadriceps', 'traps', 'triceps');"},
		{"equipment", "CREATE TYPE equipment AS ENUM ('body_only', 'dumbbell');"},
		{"difficulty_level", "CREATE TYPE difficulty_level AS ENUM ('beginner', 'intermediate', 'expert');"},
	}

	for _, enum := range enumQueries {
		if err := db.Exec(enum.Query).Error; err != nil {
			if strings.Contains(err.Error(), "already exists") {
				log.Printf("Info: Enum type '%s' already exists, skipping creation.", enum.Type)
			} else {
				log.Printf("Error: Failed to create enum type '%s': %v", enum.Type, err)
				return fmt.Errorf("failed to create enum type '%s': %w", enum.Type, err)
			}
		}
	}

	return nil
}
