package dto

import (
	"time"
)

type LogRequest struct {
	UserID     string    `json:"user_id" validate:"required"`
	ExerciseID uint      `json:"exercise_id" validate:"required_with=numeric,gt=0"`
	Timestamp  time.Time `json:"timestamp" validate:"required"`
	Set        int       `json:"set" validate:"required_with=numeric,gt=0"`
	Repetition int       `json:"repetition" validate:"required,numeric,gt=0"`
}

type WorkoutSummary struct {
	TotalWorkoutTime      time.Duration `json:"total_workout_time" validate:"required,duration"`
	TotalExercises        int           `json:"total_exercises" validate:"required,numeric"`
	TotalSets             int           `json:"total_sets" validate:"required,numeric"`
	TotalRepetitions      int           `json:"total_repetitions" validate:"required,numeric"`
	MostFrequentExercise  string        `json:"most_frequent_exercise" validate:"required"`
	AverageSessionPerWeek float64       `json:"average_session_per_week" validate:"required,numeric"`
}
