package domain

import (
	"api-peak-form/dto"
	"context"

	"gorm.io/gorm"
)

type Schedule struct {
	ID	  uint			`gorm:"primaryKey"`
    UID   string    
	Type  ExerciseType 	`gorm:"type:exercise_type"`
    Day   int          	`gorm:"not null;check:day >= 1 AND day <= 7"`
	User  User 			`gorm:"foreignKey:UID;references:ID"`

	Exercises []Exercise	`gorm:"many2many:exercise_list"`
}

type ScheduleRepository interface {
	Save(ctx context.Context, schedule *Schedule) error
	SaveExercise(ctx context.Context, schedule *ExerciseList) error
	FindByUID(ctx context.Context, ID string) ([]Schedule, error) 
	FindByUIDDayType(ctx context.Context, uid string, day int, typee string, schedule *Schedule) *Schedule
	Delete(ctx context.Context, id uint) *gorm.DB
	DeleteExercise(ctx context.Context, id uint, id_exercise int) *gorm.DB
	CountExercisesByScheduleID(ctx context.Context, id uint) int64
}

type ScheduleService interface {
	Create(ctx context.Context, req dto.CreateScheduleRequest) error
	FindByUID(ctx context.Context, uid string) (dto.ScheduleListResponse, error)
	DeleteSchedule(ctx context.Context, userID string, scheduleID uint) error
	DeleteExerciseSchedule(ctx context.Context, id uint, id_exercise int) error
}