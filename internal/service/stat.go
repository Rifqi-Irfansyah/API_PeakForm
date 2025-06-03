package service

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type statsService struct {
	logRepo domain.LogRepository
}

func NewStatService(logRepo domain.LogRepository) domain.StatsService {
	return statsService{logRepo: logRepo}
}

func (s statsService) GetStatsByUserID(ctx context.Context, userID string) (dto.StatsSummary, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logrus.Infof("Generating stats for user ID: %s", userID)

	logs, err := s.logRepo.FindByUserID(ctx, userID)
	if err != nil {
		logrus.Errorf("Failed to get logs for user ID %s: %v", userID, err)
		return dto.StatsSummary{}, fmt.Errorf("failed to get logs: %w", err)
	}

	if len(logs) == 0 {
		return dto.StatsSummary{}, fmt.Errorf("no logs found for user ID: %s", userID)
	}

	var totalSets, totalReps, totalExerciseCount int
	exerciseCount := make(map[string]int)

	for _, log := range logs {
		totalSets += log.Set
		totalReps += log.Repetition
		exerciseID := fmt.Sprintf("%d", log.ExerciseID)
		exerciseCount[exerciseID]++
		totalExerciseCount++
	}

	logrus.Infof("Stats generated: Sets=%d, Reps=%d, TotalExercises=%d", totalSets, totalReps, totalExerciseCount)

	return dto.StatsSummary{
		UserID:             userID,
		TotalSets:          totalSets,
		TotalReps:          totalReps,
		ExerciseCounter:    exerciseCount,
		TotalExerciseCount: totalExerciseCount,
	}, nil
}