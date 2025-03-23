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
	result := s.scheduleRepository.FindByUIDAndDay(ctx, req.UID, req.Day, &schedule)

	if result == nil {
		schedule = domain.Schedule{
			Type:	   domain.ExerciseType(req.Type),
			Day:       req.Day,
		}

		if err := s.scheduleRepository.Save(ctx, &schedule); err != nil {
			return err
		}

		if err := s.scheduleRepository.AddScheduleToUser(ctx, req.UID, schedule.ID); err != nil{
			return err
		}
	}

	exerciselist := domain.ExerciseList{
		ScheduleID:	schedule.ID,
		ExerciseID: req.ExerciseID,
		Set: 		uint(req.Set),
		Repetition: uint(req.Repetition),
	}

	err := s.scheduleRepository.SaveExercise(ctx, &exerciselist)
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
		for _, ex := range schedule.Exercises {
			exerciseResponses = append(exerciseResponses, dto.ExerciseResponse{
				ID:           ex.ID,
				Name:         ex.Name,
				Type:         string(ex.Type),
				Muscle:       string(ex.Muscle),
				Equipment:    string(ex.Equipment),
				Difficulty:   string(ex.Difficulty),
				Instructions: ex.Instructions,
				Gif:          ex.Gif,
			})
		}

		scheduleResponses = append(scheduleResponses, dto.ScheduleResponse{
			ID:        schedule.ID,
			Day:       schedule.Day,
			Exercises: exerciseResponses,
		})
	}

	return dto.ScheduleListResponse{
		Schedules: scheduleResponses,
	}, nil
}

func (s *scheduleService) Delete(ctx context.Context, userID string, scheduleID uint) error {
    err := s.scheduleRepository.DeleteUserSchedule(ctx, userID, scheduleID)
    if err != nil {
        return fmt.Errorf("failed to delete schedule: %w", err)
    }
    return nil
}
