package service

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

type scheduleService struct {
	scheduleRepository domain.ScheduleRepository
}

func NewScheduleService(scheduleRepository domain.ScheduleRepository) domain.ScheduleService {
	return &scheduleService{scheduleRepository: scheduleRepository}
}

func (s *scheduleService) Create(ctx context.Context, req dto.CreateScheduleRequest) error {
	logrus.Infof("Creating schedule for UID: %s, Day: %s, Type: %s", req.UID, req.Day, req.Type)
	var schedule domain.Schedule
	result := s.scheduleRepository.FindByUIDDayType(ctx, req.UID, req.Day, req.Type, &schedule)

	if result == nil {
		schedule = domain.Schedule{
			UID:  req.UID,
			Type: domain.ExerciseType(req.Type),
			Day:  req.Day,
		}

		if err := s.scheduleRepository.Save(ctx, &schedule); err != nil {
			logrus.Errorf("Failed to save schedule: %v", err)
			return err
		}
	}

	exerciselist := domain.ExerciseList{
		ScheduleID: schedule.ID,
		ExerciseID: req.ExerciseID,
		Set:        uint(req.Set),
		Repetition: uint(req.Repetition),
	}

	err := s.scheduleRepository.SaveExercise(ctx, &exerciselist)
	if err != nil {
		logrus.Errorf("Failed to save exercise list: %v", err)
		return err
	}

	logrus.Infof("Schedule created successfully for UID: %s, Day: %s, Type: %s", req.UID, req.Day, req.Type)
	return nil
}

func (s scheduleService) Update(ctx context.Context, req dto.UpdateScheduleRequest) error {
	logrus.Infof("Updating schedule with ID: %d", req.ID)
	persisted, err := s.scheduleRepository.FindById(ctx, req.ID)
	if err != nil {
		logrus.Errorf("Failed to find schedule with ID %d: %v", req.ID, err)
		return err
	}
	if persisted.ID == 0 {
		logrus.Warnf("Schedule with ID %d not found", req.ID)
		return errors.New("schedule not found")
	}
	persisted.Day = req.Day

	err = s.scheduleRepository.Update(ctx, &persisted)
	if err != nil {
		logrus.Errorf("Failed to update schedule with ID %d: %v", req.ID, err)
		return err
	}

	logrus.Infof("Schedule updated successfully with ID: %d", req.ID)
	return nil
}

func (s scheduleService) UpdateExerciseSchedule(ctx context.Context, req dto.UpdateExerciseScheduleRequest) error {
	logrus.Infof("Updating exercise schedule with Schedule ID: %d and Exercise ID: %d", req.ID, req.ExerciseID)

	persisted, err := s.scheduleRepository.FindByIdExerciseId(ctx, req.ID, req.ExerciseID)
	if err != nil {
		logrus.Errorf("Failed to find exercise schedule with Schedule ID %d and Exercise ID %d: %v", req.ID, req.ExerciseID, err)
		return err
	}
	if persisted.ScheduleID == 0 {
		logrus.Warnf("Schedule with ID %d not found", req.ID)
		return errors.New("schedule not found")
	}
	if persisted.ExerciseID == 0 {
		logrus.Warnf("Exercise with ID %d not found", req.ExerciseID)
		return errors.New("exercise not found")
	}

	updates := map[string]interface{}{}

	if req.Repetition != 0 {
		updates["repetition"] = uint(req.Repetition)
	}
	if req.Set != 0 {
		updates["set"] = uint(req.Set)
	}
	if req.NExerciseID != 0 {
		updates["exercise_id"] = uint(req.NExerciseID)
	}

	if len(updates) == 0 {
		logrus.Infof("No updates provided for Schedule ID: %d and Exercise ID: %d", req.ID, req.ExerciseID)
		return nil
	}

	err = s.scheduleRepository.UpdateExercise(ctx, persisted.ScheduleID, persisted.ExerciseID, updates)
	if err != nil {
		logrus.Errorf("Failed to update exercise schedule with Schedule ID %d and Exercise ID %d: %v", req.ID, req.ExerciseID, err)
		return err
	}

	logrus.Infof("Exercise schedule updated successfully with Schedule ID: %d and Exercise ID: %d", req.ID, req.ExerciseID)
	return nil
}

func (s scheduleService) FindByUID(ctx context.Context, uid string) (dto.ScheduleListResponse, error) {
	logrus.Infof("Fetching schedules for UID: %s", uid)
	schedules, err := s.scheduleRepository.FindByUID(ctx, uid)
	if err != nil {
		logrus.Errorf("Error saat mencari UID: %s, %v", uid, err)
		return dto.ScheduleListResponse{}, errors.New("user not found")
	}

	if len(schedules) == 0 {
		logrus.Warnf("Schedule tidak ditemukan untuk UID: %s", uid)
		return dto.ScheduleListResponse{}, errors.New("schedule not found")
	}

	logrus.Infof("Schedule ditemukan untuk user: %s", uid)

	var scheduleResponses []dto.ScheduleResponse

	for _, schedule := range schedules {
		var exerciseResponses []dto.ExerciseResponse
		for _, exList := range schedule.ExerciseList {
			exerciseResponses = append(exerciseResponses, dto.ExerciseResponse{
				ID:           exList.Exercise.ID,
				Set:          int(exList.Set),
				Repetition:   int(exList.Repetition),
				Name:         exList.Exercise.Name,
				Type:         string(exList.Exercise.Type),
				Muscle:       string(exList.Exercise.Muscle),
				Equipment:    string(exList.Exercise.Equipment),
				Difficulty:   string(exList.Exercise.Difficulty),
				Instructions: exList.Exercise.Instructions,
				Image:          exList.Exercise.Image,
			})
		}

		scheduleResponses = append(scheduleResponses, dto.ScheduleResponse{
			ID:        schedule.ID,
			Day:       schedule.Day,
			Type:      string(schedule.Type),
			Exercises: exerciseResponses,
		})
	}

	logrus.Infof("Fetched %d schedules for UID: %s", len(scheduleResponses), uid)
	return dto.ScheduleListResponse{
		Schedules: scheduleResponses,
	}, nil
}

func (s *scheduleService) DeleteSchedule(ctx context.Context, scheduleID uint) error {
	logrus.Infof("Attempting to delete schedule with ID: %d", scheduleID)
	result := s.scheduleRepository.Delete(ctx, scheduleID)
	if result.Error != nil {
		logrus.Errorf("Failed to delete schedule with ID %d: %v", scheduleID, result.Error)
		return errors.New(result.Error.Error())
	}

	if result.RowsAffected == 0 {
		logrus.Warnf("Schedule with ID %d not found", scheduleID)
		return fmt.Errorf("schedule not found")
	}

	logrus.Infof("Schedule deleted successfully with ID: %d", scheduleID)
	return nil
}

func (s *scheduleService) DeleteExerciseSchedule(ctx context.Context, id uint, id_exercise int) error {
	logrus.Infof("Attempting to delete exercise schedule with Schedule ID: %d and Exercise ID: %d", id, id_exercise)
	err := s.scheduleRepository.DeleteExercise(ctx, id, id_exercise)
	if err != nil {
		logrus.Errorf("Failed to delete exercise schedule with Schedule ID %d and Exercise ID %d: %v", id, id_exercise, err)
		return fmt.Errorf("exercise schedule not found")
	}

	count := s.scheduleRepository.CountExercisesByScheduleID(ctx, id)
	if count == 0 {
		logrus.Infof("No more exercises left in schedule with ID: %d, deleting schedule", id)
		s.scheduleRepository.Delete(ctx, id)
	}

	logrus.Infof("Exercise schedule deleted successfully with Schedule ID: %d and Exercise ID: %d", id, id_exercise)
	return nil
}
