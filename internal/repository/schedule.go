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

func (s *Schedule) FindById(ctx context.Context, id string) (result domain.Schedule, err error) {
	var schedule domain.Schedule
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&schedule).Error; err != nil {
		return domain.Schedule{}, err
	}
	return schedule, nil
}

func (s *Schedule) FindByUID(ctx context.Context, id string) ([]domain.Schedule, error) {
	var user domain.User

	err := s.db.WithContext(ctx).Preload("Schedules.Exercises").Preload("Schedules.ExerciseList").
			Preload("Schedules.ExerciseList.Exercise").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user.Schedules, nil
}

func (s *Schedule) FindByIdExerciseId(ctx context.Context, id string, eid uint) (domain.ExerciseList, error) {
	var exercise_list domain.ExerciseList
	if err := s.db.WithContext(ctx).Where("schedule_id = ? AND exercise_id = ?", id, eid).First(&exercise_list).Error; err != nil {
		return domain.ExerciseList{}, err
	}
	return exercise_list, nil
}

func (sc *Schedule) FindByUIDDayType(ctx context.Context, uid string, day int, typee string, schedule *domain.Schedule) *domain.Schedule {
	if err := sc.db.WithContext(ctx).Where("uid = ? AND day = ? AND type = ?", uid, day, typee).First(schedule).Error; err != nil {
        return nil
    }

	return schedule
}

func (sc *Schedule) Save(ctx context.Context, s *domain.Schedule) error {
	return sc.db.WithContext(ctx).Create(s).Error
}

func (sc *Schedule) SaveExercise(ctx context.Context, s *domain.ExerciseList) error {
	return sc.db.WithContext(ctx).Create(s).Error
}

func (sc *Schedule) Update(ctx context.Context, s *domain.Schedule) error {
	return sc.db.WithContext(ctx).Save(s).Error
}

func (sc *Schedule) UpdateExercise(ctx context.Context, scheduleID, exerciseID uint, updates map[string]interface{}) error {
    return sc.db.WithContext(ctx).
        Model(&domain.ExerciseList{}).
        Where("schedule_id = ? AND exercise_id = ?", scheduleID, exerciseID).
        Updates(updates).
        Error
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