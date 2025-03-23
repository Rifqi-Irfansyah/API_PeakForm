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

	err := r.db.WithContext(ctx).
		Preload("Schedules.Exercises"). // Preload relasi ke Exercises juga (jika ada)
		Where("id = ?", ID).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return user.Schedules, nil
}

func (sc *Schedule) FindByUIDDayType(ctx context.Context, uid string, day int, typee string, schedule *domain.Schedule) *domain.Schedule {
	var idSchedules []uint

	if err := sc.db.WithContext(ctx).Table("user_schedules").Where("user_id = ?", uid).Pluck("schedule_id", &idSchedules).Error; err != nil {
		fmt.Errorf("failed to fetch schedule IDs: %w", err)
        return nil
    }

    // Jika tidak ada schedule_id yang ditemukan, kembalikan error
    if len(idSchedules) == 0 {
		fmt.Errorf("no schedules found for user %s", uid)
        return nil
    }

    // Query ke tabel schedules berdasarkan schedule_id, day, dan type
    if err := sc.db.WithContext(ctx).Where("id IN (?) AND day = ? AND type = ?", idSchedules, day, typee).First(schedule).Error; err != nil {
        fmt.Errorf("failed to fetch schedule: %w", err)
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
	var count int64
	sc.db.WithContext(ctx).Table("user_schedules").Where("schedule_id = ?", id).Count(&count)

	if count > 0 {
		return &gorm.DB{}
	}
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
        return result.Error
    }

    return nil
}