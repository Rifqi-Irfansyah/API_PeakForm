package migration

import (
	"api-peak-form/domain"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Migrate(dbConnection *gorm.DB) error {
	logrus.Info("Starting database migration")

	dbConnection.Exec("DROP TABLE IF EXISTS exercise_list, user_schedules, users, schedules, exercises, logs, exercise_lists CASCADE;")
	logrus.Info("Dropped existing tables if they exist")

	if err := createEnums(dbConnection); err != nil {
		logrus.Warnf("Warning: Failed to create enums: %v", err)
	}

	err := dbConnection.AutoMigrate(&domain.User{}, &domain.Log{}, &domain.Exercise{}, &domain.Schedule{}, &domain.ExerciseList{})
	if err != nil {
		logrus.Errorf("Failed to migrate database: %v", err)
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	logrus.Info("Database migrated successfully")
	return nil
}

func createEnums(db *gorm.DB) error {
	logrus.Info("Starting to create enums")
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
		logrus.Infof("Creating enum type '%s'", enum.Type)
		if err := db.Exec(enum.Query).Error; err != nil {
			logrus.Errorf("Failed to create enum type '%s': %v", enum.Type, err)
			return fmt.Errorf("failed to create enum type '%s': %w", enum.Type, err)
		}
		logrus.Infof("Successfully created enum type '%s'", enum.Type)
	}

	logrus.Info("Finished creating enums")
	return nil
}
