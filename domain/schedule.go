package domain

import (
	"api-peak-form/dto"
	"context"

	"gorm.io/gorm"
)

type Schedule struct {
	ID        uint 		 	`gorm:"primaryKey"`
	Type	  ExerciseType	`gorm:"type:exercise_type"`
	Day       int         	`gorm:"not null; check:day >= 1 AND day <= 7"`

	Exercises []Exercise	`gorm:"many2many:exercise_list"`
	User	  []User		`gorm:"many2many:user_schedules"`
}

type ScheduleRepository interface {
	Save(ctx context.Context, schedule *Schedule) error
	SaveExercise(ctx context.Context, schedule *ExerciseList) error
	AddScheduleToUser(ctx context.Context, userID string, scheduleID uint) error
	FindByUID(ctx context.Context, ID string) ([]Schedule, error) 
	FindByUIDAndDay(ctx context.Context, uid string, day int, schedule *Schedule) *Schedule
	Delete(ctx context.Context, id uint) *gorm.DB
	DeleteExercise(ctx context.Context, id uint) *gorm.DB
	DeleteUserSchedule(ctx context.Context, userID string, scheduleID uint) error

}

type ScheduleService interface {
	// Save(ctx context.Context, schedule *Schedule) error
	Create(ctx context.Context, req dto.CreateScheduleRequest) error
	FindByUID(ctx context.Context, uid string) (dto.ScheduleListResponse, error)
	Delete(ctx context.Context, userID string, scheduleID uint) error
}