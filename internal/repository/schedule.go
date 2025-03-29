package repository

import (
	"api-peak-form/domain"
	"context"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type Schedule struct {
	db *gorm.DB
}

func NewSchedule(db *gorm.DB) *Schedule {
	logrus.Info("Creating new Schedule")
	return &Schedule{db: db}
}

func (s *Schedule) FindById(ctx context.Context, id string) (result domain.Schedule, err error) {
	logrus.Infof("Finding schedule by ID: %s", id)
	var schedule domain.Schedule
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&schedule).Error; err != nil {
		logrus.Errorf("Error finding schedule by ID %s: %v", id, err)
		return domain.Schedule{}, err
	}
	logrus.Infof("Schedule found successfully by ID: %s", id)
	return schedule, nil
}

func (s *Schedule) FindByUID(ctx context.Context, id string) ([]domain.Schedule, error) {
	logrus.Infof("Finding schedules by user ID: %s", id)
	var user domain.User

	err := s.db.WithContext(ctx).Preload("Schedules.Exercises").Preload("Schedules.ExerciseList").
		Preload("Schedules.ExerciseList.Exercise").Where("id = ?", id).First(&user).Error
	if err != nil {
		logrus.Errorf("Error finding schedules for user ID %s: %v", id, err)
		return nil, err
	}

	logrus.Infof("Schedules found successfully for user ID: %s", id)
	return user.Schedules, nil
}

func (s *Schedule) FindByIdExerciseId(ctx context.Context, id string, eid uint) (domain.ExerciseList, error) {
	logrus.Infof("Finding exercise list by schedule ID: %s and exercise ID: %d", id, eid)
	var exerciseList domain.ExerciseList
	if err := s.db.WithContext(ctx).Where("schedule_id = ? AND exercise_id = ?", id, eid).First(&exerciseList).Error; err != nil {
		logrus.Errorf("Error finding exercise list by schedule ID %s and exercise ID %d: %v", id, eid, err)
		return domain.ExerciseList{}, err
	}
	logrus.Infof("Exercise list found successfully by schedule ID: %s and exercise ID: %d", id, eid)
	return exerciseList, nil
}

func (sc *Schedule) FindByUIDDayType(ctx context.Context, uid string, day int, typee string, schedule *domain.Schedule) *domain.Schedule {
	logrus.Infof("Finding schedule by UID: %s, day: %d, type: %s", uid, day, typee)
	if err := sc.db.WithContext(ctx).Where("uid = ? AND day = ? AND type = ?", uid, day, typee).First(schedule).Error; err != nil {
		logrus.Errorf("Error finding schedule by UID %s, day %d, type %s: %v", uid, day, typee, err)
		return nil
	}
	logrus.Infof("Schedule found successfully by UID: %s, day: %d, type: %s", uid, day, typee)
	return schedule
}

func (sc *Schedule) Save(ctx context.Context, s *domain.Schedule) error {
	logrus.Infof("Saving schedule: %v", s)
	err := sc.db.WithContext(ctx).Create(s).Error
	if err != nil {
		logrus.Errorf("Error saving schedule: %v", err)
		return err
	}
	logrus.Infof("Schedule saved successfully: %v", s)
	return nil
}

func (sc *Schedule) SaveExercise(ctx context.Context, s *domain.ExerciseList) error {
	logrus.Infof("Saving exercise: %v", s)
	err := sc.db.WithContext(ctx).Create(s).Error
	if err != nil {
		logrus.Errorf("Error saving exercise: %v", err)
		return err
	}
	logrus.Infof("Exercise saved successfully: %v", s)
	return nil
}

func (sc *Schedule) Update(ctx context.Context, s *domain.Schedule) error {
	logrus.Infof("Updating schedule: %v", s)
	err := sc.db.WithContext(ctx).Save(s).Error
	if err != nil {
		logrus.Errorf("Error updating schedule: %v", err)
		return err
	}
	logrus.Infof("Schedule updated successfully: %v", s)
	return nil
}

func (sc *Schedule) UpdateExercise(ctx context.Context, scheduleID, exerciseID uint, updates map[string]interface{}) error {
	logrus.Infof("Updating exercise with schedule ID: %d and exercise ID: %d", scheduleID, exerciseID)
	err := sc.db.WithContext(ctx).
		Model(&domain.ExerciseList{}).
		Where("schedule_id = ? AND exercise_id = ?", scheduleID, exerciseID).
		Updates(updates).
		Error
	if err != nil {
		logrus.Errorf("Error updating exercise with schedule ID: %d and exercise ID: %d: %v", scheduleID, exerciseID, err)
		return err
	}
	logrus.Infof("Exercise updated successfully with schedule ID: %d and exercise ID: %d", scheduleID, exerciseID)
	return nil
}

func (sc *Schedule) Delete(ctx context.Context, id uint) *gorm.DB {
	logrus.Infof("Deleting schedule with ID: %d", id)
	result := sc.db.WithContext(ctx).Delete(&domain.Schedule{}, id)
	if result.Error != nil {
		logrus.Errorf("Error deleting schedule with ID %d: %v", id, result.Error)
	} else {
		logrus.Infof("Schedule deleted successfully with ID: %d", id)
	}
	return result
}

func (sc *Schedule) DeleteExercise(ctx context.Context, id uint, id_exercise int) *gorm.DB {
	logrus.Infof("Deleting exercise with schedule ID: %d and exercise ID: %d", id, id_exercise)
	result := sc.db.WithContext(ctx).Table("exercise_list").Where("schedule_id = ? AND exercise_id = ?", id, id_exercise).Delete(nil)

	if result.Error != nil {
		logrus.Errorf("Error deleting exercise with schedule ID %d and exercise ID %d: %v", id, id_exercise, result.Error)
		return result
	}

	if result.RowsAffected == 0 {
		logrus.Warnf("No exercise found to delete with schedule ID %d and exercise ID %d", id, id_exercise)
		return result
	}

	logrus.Infof("Exercise deleted successfully with schedule ID: %d and exercise ID: %d", id, id_exercise)
	return nil
}

func (r *Schedule) CountExercisesByScheduleID(ctx context.Context, id uint) int64 {
	logrus.Infof("Counting exercises for schedule ID: %d", id)
	var count int64
	result := r.db.WithContext(ctx).
		Table("exercise_list").
		Where("schedule_id = ?", id).
		Count(&count)

	if result.Error != nil {
		logrus.Errorf("Error counting exercises for schedule ID %d: %v", id, result.Error)
		return 0
	}
	logrus.Infof("Counted %d exercises for schedule ID: %d", count, id)
	return count
}
