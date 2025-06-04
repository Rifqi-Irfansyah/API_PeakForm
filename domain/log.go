package domain

import (
	"api-peak-form/dto"
	"context"
	"time"
)

type Log struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     string    `gorm:"type:uuid;not null"`
	ExerciseID uint      `gorm:"not null"`
	Timestamp  time.Time `gorm:"not null"`
	Exercise   Exercise  `gorm:"foreignKey:ExerciseID"`
	Set        int       `gorm:"not null"`
	Repetition int       `gorm:"not null"`
}

type LogRepository interface {
	Create(ctx context.Context, log Log) error
	FindByUserID(ctx context.Context, userID string) ([]Log, error)
	FindLastByUserID(ctx context.Context, userID string) (Log, error)
}

type LogService interface {
	Create(ctx context.Context, req dto.LogRequest) error
	FindByUserID(ctx context.Context, userID string) ([]Log, error)
	GetUserWorkoutSummary(ctx context.Context, userID string) (dto.WorkoutSummary, error)
	HasUserExercisedToday(ctx context.Context, userID string) (bool, error)
}
