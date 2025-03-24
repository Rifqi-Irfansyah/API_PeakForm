package repository

import (
	"api-peak-form/domain"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Schedule struct {
	db *gorm.DB
}

func NewSchedule(db *gorm.DB) *Schedule {
	return &Schedule{db: db}
}

func (r *Schedule) FindByUID(ctx context.Context, ID string) ([]domain.Schedule, error) {
	var user domain.User

	err := r.db.WithContext(ctx).Preload("Schedules.Exercises").Where("id = ?", ID).First(&user).Error

	if err != nil {
		return nil, err
	}

	return user.Schedules, nil
}

func (sc *Schedule) FindByUIDDayType(ctx context.Context, uid string, day int, typee string, schedule *domain.Schedule) *domain.Schedule {
	if err := sc.db.WithContext(ctx).Where("uid = ? AND day = ? AND type = ?", uid, day, typee).First(schedule).Error; err != nil {
        return nil
    }

	return schedule
}

func (sc *Schedule) Save(ctx context.Context, c *domain.Schedule) error {
	return sc.db.WithContext(ctx).Create(c).Error
}

func (sc *Schedule) SaveExercise(ctx context.Context, c *domain.ExerciseList) error {
	return sc.db.WithContext(ctx).Create(c).Error
}


func (sc *Schedule) Delete(ctx context.Context, id uint) *gorm.DB {
	return sc.db.WithContext(ctx).Delete(&domain.Schedule{}, id)
}

func (sc *Schedule) DeleteExercise(ctx context.Context, id uint, id_exercise int) *gorm.DB {
	result := sc.db.WithContext(ctx).Table("exercise_list").Where("schedule_id = ? AND exercise_id = ?", id, id_exercise).Delete(nil)

	if result.Error != nil {
        fmt.Printf("Gagal: %v\n", result.Error)
        return result
    }

    if result.RowsAffected == 0 {
		return result
    }

	return nil
}

func (r *Schedule) CountExercisesByScheduleID(ctx context.Context, id uint) int64 {
    var count int64
    result := r.db.WithContext(ctx).
        Table("exercise_list").
        Where("schedule_id = ?", id).
        Count(&count)

    if result.Error != nil {
        return 0
    }
    return count
}