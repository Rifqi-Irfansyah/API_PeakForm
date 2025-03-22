package domain

import (
	"api-peak-form/dto"
	"context"
)

type Schedule struct {
	ID        uint  `gorm:"primaryKey"`
	UserID    string
	User      User       `gorm:"foreignKey:UserID"`
	Day       int        `gorm:"not null; check:day >= 1 AND day <= 7"`
	ExerciseID uint
	Exercises []Exercise `gorm:"many2many:exercise_list"`
}

type ScheduleRepository interface {
	// AddSchedule(ctx context.Context, UID string ) ()
	Save(ctx context.Context, schedule *Schedule) error
	SaveExercise(ctx context.Context, schedule *ExerciseList) error
	FindByUID(ctx context.Context, ID string) ([]Schedule, error) 
	FindByUIDAndDay(ctx context.Context, uid string, day int, schedule *Schedule) *Schedule
}

type ScheduleService interface {
	// Save(ctx context.Context, schedule *Schedule) error
	Create(ctx context.Context, req dto.CreateScheduleRequest) error
	FindByUID(ctx context.Context, uid string) (dto.ScheduleListResponse, error)
}