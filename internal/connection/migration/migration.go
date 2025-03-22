package migration

import (
	"api-peak-form/domain"
	"fmt"
	"gorm.io/gorm"
	"log"
)

func Migrate(dbConnection *gorm.DB) error {
	err := dbConnection.Migrator().DropTable(&domain.User{}, &domain.Log{}, &domain.Exercise{}, &domain.Schedule{}, &domain.ExerciseList{})
	if err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	if err := createEnums(dbConnection); err != nil {
		log.Printf("Warning: Failed to create enums: %v", err)
	}

	err = dbConnection.AutoMigrate(&domain.User{}, &domain.Log{}, &domain.Exercise{}, &domain.Schedule{}, &domain.ExerciseList{}, &domain.UserSchedule{})
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
		{"exercise_type", "DO $$ BEGIN CREATE TYPE exercise_type AS ENUM ('strength', 'cardio'); EXCEPTION WHEN duplicate_object THEN null; END $$;"},
		{"muscle_group", "DO $$ BEGIN CREATE TYPE muscle_group AS ENUM ('abdominals', 'biceps', 'calves', 'chest', 'forearms', 'lats', 'lower_back', 'middle_back', 'neck', 'quadriceps', 'traps', 'triceps'); EXCEPTION WHEN duplicate_object THEN null; END $$;"},
		{"equipment", "DO $$ BEGIN CREATE TYPE equipment AS ENUM ('body_only', 'dumbbell'); EXCEPTION WHEN duplicate_object THEN null; END $$;"},
		{"difficulty_level", "DO $$ BEGIN CREATE TYPE difficulty_level AS ENUM ('beginner', 'intermediate', 'expert'); EXCEPTION WHEN duplicate_object THEN null; END $$;"},
	}

	for _, enum := range enumQueries {
		if err := db.Exec(enum.Query).Error; err != nil {
			log.Printf("Error: Failed to create enum type '%s': %v", enum.Type, err)
			return fmt.Errorf("failed to create enum type '%s': %w", enum.Type, err)
		}
	}

	return nil
}