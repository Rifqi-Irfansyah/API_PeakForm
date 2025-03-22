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

func (r *Schedule) FindByUIDAndDay(ctx context.Context, uid string, day int, schedule *domain.Schedule) *domain.Schedule {
	if err := r.db.WithContext(ctx).Where("user_id = ? AND day = ?", uid, day).First(schedule).Error; err != nil {
		return nil 
	}
	return schedule
}


func (cr *Schedule) Save(ctx context.Context, c *domain.Schedule) error {
	return cr.db.WithContext(ctx).Create(c).Error
}

func (cr *Schedule) SaveExercise(ctx context.Context, c *domain.ExerciseList) error {
	return cr.db.WithContext(ctx).Create(c).Error
}