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
	Exercises []Exercise `gorm:"many2many:exercise_list"`
}

type ScheduleRepository interface {
	FindByUID(ctx context.Context, ID string) ([]Schedule, error) 
}

type ScheduleService interface {
	FindByUID(ctx context.Context, uid string) (dto.ScheduleListResponse, error)
}