package domain

import (
	"api-peak-form/dto"
	"context"
)

type Exercise struct {
	ID           uint            `gorm:"primaryKey"`
	Name         string          `gorm:"type:varchar(100);not null"`
	Type         ExerciseType    `gorm:"type:exercise_type;not null"`
	Muscle       MuscleGroup     `gorm:"type:muscle_group;not null"`
	Equipment    Equipment       `gorm:"type:equipment;not null"`
	Difficulty   DifficultyLevel `gorm:"type:difficulty_level;not null"`
	Instructions string          `gorm:"type:TEXT;not null"`
	Gif          string          `gorm:"type:varchar(255);not null"`
}

type ExerciseRepository interface {
	Create(ctx context.Context, exercise *Exercise) error
	GetAll(ctx context.Context) ([]Exercise, error)
	GetByID(ctx context.Context, id uint) (*Exercise, error)
	Update(ctx context.Context, exercise *Exercise) error
	Delete(ctx context.Context, id uint) error
}

type ExerciseService interface {
	CreateExercise(ctx context.Context, req dto.CreateExerciseRequest) error
	GetExercises(ctx context.Context) ([]Exercise, error)
	GetExerciseByID(ctx context.Context, id uint) (Exercise, error)
	UpdateExercise(ctx context.Context, req dto.UpdateExerciseRequest) error
	DeleteExercise(ctx context.Context, id uint) error
}