package service

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"context"
	"fmt"
	"time"
)

type logService struct {
	logRepository  domain.LogRepository
	userRepository domain.UserRepository
}

func NewLogService(logRepository domain.LogRepository, userRepository domain.UserRepository) domain.LogService {
	return logService{
		logRepository:  logRepository,
		userRepository: userRepository,
	}
}

func (l logService) Create(ctx context.Context, req dto.LogRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	log := domain.Log{
		UserID:     req.UserID,
		ExerciseID: req.ExerciseID,
		Timestamp:  req.Timestamp,
		Set:        req.Set,
		Repetition: req.Repetition,
	}

	_, err := l.userRepository.FindByID(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("failed to find user with ID %s: %w", req.UserID, err)
	}

	err = l.logRepository.Create(ctx, log)
	if err != nil {
		return fmt.Errorf("failed to create log: %w", err)
	}

	return nil
}

func (l logService) FindByUserID(ctx context.Context, userID string) ([]domain.Log, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logs, err := l.logRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch logs for user ID %d: %w", userID, err)
	}

	return logs, nil
}
