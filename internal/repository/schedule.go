package repository

import (
	"api-peak-form/domain"
	"context"

	"gorm.io/gorm"
)

type Schedule struct {
	db *gorm.DB
}

func NewSchedule(db *gorm.DB) *Schedule {
	return &Schedule{db: db}
}

func (r *Schedule) FindByUID(ctx context.Context, ID string) ([]domain.Schedule, error) {
	var schedules []domain.Schedule
	err := r.db.WithContext(ctx).Where("user_id = ?", ID).Preload("Exercises").Find(&schedules).Error
	if err != nil {
		return schedules, err
	}
	return schedules, nil
}