package service

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"context"
	"fmt"
	"time"
)

type logService struct {
	logRepository domain.LogRepository
}

func NewLogService(logRepository domain.LogRepository) domain.LogService {
	return logService{logRepository: logRepository}
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

	err := l.logRepository.Create(ctx, log)
	if err != nil {
		return fmt.Errorf("failed to create log: %w", err)
	}

	return nil
}

func (l logService) FindByUserID(ctx context.Context, userID uint) ([]domain.Log, error) {
	//TODO implement me
	panic("implement me")
}
