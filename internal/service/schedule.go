package service

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"context"
	"errors"
	"fmt"
	"log"
)

type scheduleService struct {
	scheduleRepository domain.ScheduleRepository
}

func NewScheduleService(scheduleRepository domain.ScheduleRepository) domain.ScheduleService {
	return &scheduleService{scheduleRepository: scheduleRepository}
}

func (s *scheduleService) Create(ctx context.Context, req dto.CreateScheduleRequest) error {
	var schedule domain.Schedule
	result := s.scheduleRepository.FindByUIDDayType(ctx, req.UID, req.Day, req.Type, &schedule)

	if result == nil {
		schedule = domain.Schedule{
			UID:  req.UID,
			Type: domain.ExerciseType(req.Type),
			Day:  req.Day,
		}

		if err := s.scheduleRepository.Save(ctx, &schedule); err != nil {
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
		return err
	}

	return nil
}

func (s scheduleService) Update(ctx context.Context, req dto.UpdateScheduleRequest) error {
	persisted, err := s.scheduleRepository.FindById(ctx, req.ID)
	if err != nil {
		return err
	}
	if persisted.ID == 0 {
		return errors.New("schedule not found")
	}
	persisted.Day = req.Day

	err = s.scheduleRepository.Update(ctx, &persisted)
	if err != nil {
		return err
	}
	return nil
}

func (s scheduleService) UpdateExerciseSchedule(ctx context.Context, req dto.UpdateExerciseScheduleRequest) error {
	persisted, err := s.scheduleRepository.FindByIdExerciseId(ctx, req.ID, req.ExerciseID)
	if err != nil {
		return err
	}
	if persisted.ScheduleID == 0 {
		return errors.New("schedule not found")
	}
	if persisted.ExerciseID == 0 {
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
        return nil
    }

	err = s.scheduleRepository.UpdateExercise(ctx, persisted.ScheduleID, persisted.ExerciseID, updates)
	if err != nil {
		return err
	}
	return nil
}

func (s scheduleService) FindByUID(ctx context.Context, uid string) (dto.ScheduleListResponse, error) {
	schedules, err := s.scheduleRepository.FindByUID(ctx, uid)
	if err != nil {
		log.Println("Error saat mencari UID: ", err)
		return dto.ScheduleListResponse{}, errors.New("user not found")
	}

	if len(schedules) == 0 {
		log.Println("Schedule tidak ditemukan:", uid)
		return dto.ScheduleListResponse{}, errors.New("schedule not found")
	}

	log.Println("Schedule ditemukan untuk user:", uid)

	var scheduleResponses []dto.ScheduleResponse

	for _, schedule := range schedules {
		var exerciseResponses []dto.ExerciseResponse
		for _, exList := range schedule.ExerciseList {
			exerciseResponses = append(exerciseResponses, dto.ExerciseResponse{
				ID:				exList.Exercise.ID,
				Set:			int(exList.Set),
				Repetition:		int(exList.Repetition),
				Name:			exList.Exercise.Name,
				Type:			string(exList.Exercise.Type),
				Muscle:			string(exList.Exercise.Muscle),
				Equipment:		string(exList.Exercise.Equipment),
				Difficulty:		string(exList.Exercise.Difficulty),
				Instructions:	exList.Exercise.Instructions,
				Gif:			exList.Exercise.Gif,
			})
		}

		scheduleResponses = append(scheduleResponses, dto.ScheduleResponse{
			ID:        schedule.ID,
			Day:       schedule.Day,
			Type:      string(schedule.Type),
			Exercises: exerciseResponses,
		})
	}

	return dto.ScheduleListResponse{
		Schedules: scheduleResponses,
	}, nil
}

func (s *scheduleService) DeleteSchedule(ctx context.Context, userID string, scheduleID uint) error {
	result := s.scheduleRepository.Delete(ctx, scheduleID)
	if result.Error != nil {
		return errors.New(result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("schedule not found")
	}

	return nil
}

func (s *scheduleService) DeleteExerciseSchedule(ctx context.Context, id uint, id_exercise int) error {
	err := s.scheduleRepository.DeleteExercise(ctx, id, id_exercise)
	if err != nil {
		return fmt.Errorf("exercise schedule not found")
	}

	count := s.scheduleRepository.CountExercisesByScheduleID(ctx, id)
	if count == 0 {
		s.scheduleRepository.Delete(ctx, id)
	}

	return nil
}
