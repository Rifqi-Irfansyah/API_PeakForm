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

type User struct {
	db *gorm.DB
}

func NewSchedule(db *gorm.DB) *Schedule {
	return &Schedule{db: db}
}

func (r *Schedule) FindByUID(ctx context.Context, ID string) ([]domain.Schedule, error) {
	var user domain.User

	err := r.db.WithContext(ctx).
		Preload("Schedules.Exercises"). // Preload relasi ke Exercises juga (jika ada)
		Where("id = ?", ID).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return user.Schedules, nil
}

func (sc *Schedule) FindByUIDAndDay(ctx context.Context, uid string, day int, schedule *domain.Schedule) *domain.Schedule {
	if err := sc.db.WithContext(ctx).Where("user_id = ? AND day = ?", uid, day).First(schedule).Error; err != nil {
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

func (sc *Schedule) AddScheduleToUser(ctx context.Context, userID string, scheduleID uint) error {
	type UserSchedule struct {
		UserID     string `gorm:"column:user_id"`
		ScheduleID uint   `gorm:"column:schedule_id"`
	}
	
	userschedule := UserSchedule{
		UserID:     userID,
		ScheduleID: scheduleID,
	}
	
	err := sc.db.Table("user_schedules").Create(&userschedule).Error
	if err != nil {
		return err
	}

	return nil
}

func (sc *Schedule) Delete(ctx context.Context, id uint) *gorm.DB {
	return sc.db.WithContext(ctx).Delete(&domain.Schedule{}, id)
}

func (sc *Schedule) DeleteExercise(ctx context.Context, id uint) *gorm.DB {
	return sc.db.WithContext(ctx).Delete(&domain.ExerciseList{}, id)
}

func (sc *Schedule) DeleteUserSchedule(ctx context.Context, userID string, scheduleID uint) error {
	type UserSchedule struct {
		UserID     string `gorm:"column:user_id"`
		ScheduleID uint   `gorm:"column:schedule_id"`
	}
	
    userSchedule := UserSchedule{
        UserID:    userID,
        ScheduleID: scheduleID,
    }

    result := sc.db.WithContext(ctx).Where("user_id = ? AND schedule_id = ?", userID, scheduleID).Delete(&userSchedule)
    if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return fmt.Errorf("no record found for user_id %s and schedule_id %d", userID, scheduleID)
    }

    return nil
}