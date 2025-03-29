package service

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
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

	logrus.Infof("Creating log for user ID: %s and exercise ID: %d", req.UserID, req.ExerciseID)

	log := domain.Log{
		UserID:     req.UserID,
		ExerciseID: req.ExerciseID,
		Timestamp:  req.Timestamp,
		Set:        req.Set,
		Repetition: req.Repetition,
	}

	_, err := l.userRepository.FindByID(ctx, req.UserID)
	if err != nil {
		logrus.Errorf("Failed to find user with ID %s: %v", req.UserID, err)
		return fmt.Errorf("failed to find user with ID %s: %w", req.UserID, err)
	}

	err = l.logRepository.Create(ctx, log)
	if err != nil {
		logrus.Errorf("Failed to create log: %v", err)
		return fmt.Errorf("failed to create log: %w", err)
	}

	logrus.Infof("Log created successfully for user ID: %s and exercise ID: %d", req.UserID, req.ExerciseID)
	return nil
}

func (l logService) FindByUserID(ctx context.Context, userID string) ([]domain.Log, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logrus.Infof("Fetching logs for user ID: %s", userID)
	logs, err := l.logRepository.FindByUserID(ctx, userID)
	if err != nil {
		logrus.Errorf("Failed to fetch logs for user ID %s: %v", userID, err)
		return nil, fmt.Errorf("failed to fetch logs for user ID %s: %w", userID, err)
	}

	logrus.Infof("Fetched %d logs for user ID: %s", len(logs), userID)
	return logs, nil
}

func (l logService) GetUserWorkoutSummary(ctx context.Context, userID string) (dto.WorkoutSummary, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logrus.Infof("Fetching workout summary for user ID: %s", userID)

	logs, err := l.logRepository.FindByUserID(ctx, userID)
	if err != nil {
		logrus.Errorf("Failed to fetch logs for user ID %s: %v", userID, err)
		return dto.WorkoutSummary{}, fmt.Errorf("failed to fetch logs for user ID %s: %w", userID, err)
	}

	var totalWorkoutTime time.Duration
	var totalExercises, totalSets, totalRepetitions int
	exerciseCount := make(map[string]int)

	for _, log := range logs {
		totalWorkoutTime += log.Timestamp.Sub(log.Timestamp)
		totalExercises++
		totalSets += log.Set
		totalRepetitions += log.Repetition
		exerciseCount[log.Exercise.Name]++
	}

	var mostFrequentExercise string
	var maxCount int
	for exercise, count := range exerciseCount {
		if count > maxCount {
			maxCount = count
			mostFrequentExercise = exercise
		}
	}

	averageSessionPerWeek := float64(totalExercises) / 7.0

	logrus.Infof("Workout summary for user ID %s: TotalWorkoutTime=%v, TotalExercises=%d, TotalSets=%d, TotalRepetitions=%d, MostFrequentExercise=%s, AverageSessionPerWeek=%.2f",
		userID, totalWorkoutTime, totalExercises, totalSets, totalRepetitions, mostFrequentExercise, averageSessionPerWeek)

	return dto.WorkoutSummary{
		TotalWorkoutTime:      totalWorkoutTime,
		TotalExercises:        totalExercises,
		TotalSets:             totalSets,
		TotalRepetitions:      totalRepetitions,
		MostFrequentExercise:  mostFrequentExercise,
		AverageSessionPerWeek: averageSessionPerWeek,
	}, nil
}
