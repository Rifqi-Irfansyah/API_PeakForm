package service

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type userService struct {
	userRepository     domain.UserRepository
	scheduleRepository domain.ScheduleRepository
	logRepository      domain.LogRepository
}

func NewUserService(userRepository domain.UserRepository, scheduleRepository domain.ScheduleRepository, logRepository domain.LogRepository) domain.UserService {
	return userService{
		userRepository:     userRepository,
		scheduleRepository: scheduleRepository,
		logRepository:      logRepository,
	}
}

func (u userService) UpdatePoint(ctx context.Context, id string, difficulty domain.DifficultyLevel, rep int, set int) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	point := calculatePoints(difficulty, rep, set)
	logrus.Infof("Total point for user ID: %s is %d", id, point)

	existingUser, err := u.userRepository.FindByID(ctx, id)
	if err != nil {
		logrus.Errorf("Failed to find user with ID %s: %v", id, err)
		return 0, fmt.Errorf("failed to find user with ID %s: %w", id, err)
	}

	existingUser.Point = calculateTotalPoints(existingUser.Point, point, existingUser.Streak)
	logrus.Infof("User %s point updated to %d", id, existingUser.Point)

	if err := u.userRepository.UpdatePoint(ctx, id, existingUser.Point); err != nil {
		logrus.Errorf("Failed to update point for user ID %s: %v", id, err)
		return 0, fmt.Errorf("failed to update point for user ID %s: %w", id, err)
	}

	logrus.Infof("Point updated successfully for user ID: %s", id)
	return point, nil
}

func (u userService) CheckStreak(ctx context.Context, id string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	existingUser, err := u.userRepository.FindByID(ctx, id)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to find user with ID %s", id)
		return 0, fmt.Errorf("failed to find user with ID %s: %w", id, err)
	}

	schedules, err := u.scheduleRepository.FindByUID(ctx, id)
	if err != nil {
		logrus.Errorf("Failed to find schedules for user ID %s: %v", id, err)
		return 0, fmt.Errorf("failed to find schedules for user ID %s: %w", id, err)
	}
	if len(schedules) == 0 {
		logrus.Warnf("No schedules found for user ID %s", id)
		return 0, fmt.Errorf("no schedules found for user ID %s", id)
	}

	lastLog, err := u.logRepository.FindLastByUserID(ctx, id)
	if err != nil {
		logrus.Errorf("Failed to find last log for user ID %s: %v", id, err)
		return 0, fmt.Errorf("failed to find last log for user ID %s: %w", id, err)
	}

	if lastLog.ID == 0 {
		logrus.Infof("User %s has no log history, skipping streak update logic", id)
		existingUser.Streak = 0
	}

	lastDate := lastLog.Timestamp.Truncate(24 * time.Hour)
	today := time.Now().Truncate(24 * time.Hour)
	days := int(today.Sub(lastDate).Hours() / 24)

	reset := false

	for i := 1; i < days; i++ {
		checkDate := lastDate.AddDate(0, 0, i)
		checkDay := convertWeekdayToInt(checkDate.Weekday())
		if checkScheduleOnDay(schedules, checkDay) {
			logrus.Infof("User %s missed scheduled day: %v", id, checkDate)
			reset = true
			break
		}
	}

	if reset {
		existingUser.Streak = 0
		logrus.Infof("User %s missed required exercise, streak reset", id)
	} else {
		logrus.Infof("User %s did not miss required exercise, streak maintained", id)
	}

	if err := u.userRepository.UpdateStreak(ctx, id, existingUser.Streak); err != nil {
		logrus.WithError(err).Errorf("Failed to update streak for user ID %s", id)
		return 0, fmt.Errorf("failed to update streak for user ID %s: %w", id, err)
	}

	logrus.Infof("Streak updated successfully for user ID: %s", id)
	return existingUser.Streak, nil
}

func (u userService) UpdateStreak(ctx context.Context, id string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	existingUser, err := u.userRepository.FindByID(ctx, id)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to find user with ID %s", id)
		return 0, fmt.Errorf("failed to find user with ID %s: %w", id, err)
	}

	if existingUser.Streak < 10 {
		existingUser.Streak++
	}

	if err := u.userRepository.UpdateStreak(ctx, id, existingUser.Streak); err != nil {
		logrus.WithError(err).Errorf("Failed to increment streak for user ID %s", id)
		return 0, fmt.Errorf("failed to increment streak for user ID %s: %w", id, err)
	}

	logrus.Infof("Streak incremented successfully for user ID: %s", id)
	return existingUser.Streak, nil
}

func convertWeekdayToInt(day time.Weekday) int {
	if day == time.Sunday {
		return 7
	}
	return int(day)
}

func checkScheduleOnDay(schedules []domain.Schedule, day int) bool {
	for _, schedule := range schedules {
		if schedule.Day == day {
			return true
		}
	}
	return false
}

func calculatePoints(difficulty domain.DifficultyLevel, rep, set int) int {
	basePoint := map[domain.DifficultyLevel]int{
		domain.Beginner:     1,
		domain.Intermediate: 2,
		domain.Expert:       3,
	}[difficulty]
	return basePoint * rep * set
}

func calculateTotalPoints(currentPoints, newPoints, streak int) int {
	if streak == 0 {
		return currentPoints + newPoints
	}
	return currentPoints + (newPoints * streak)
}

func (u userService) GetAllUsersDesc(ctx context.Context) ([]dto.UserLeaderboardResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	users, err := u.userRepository.GetAllUsersDesc(ctx)
	if err != nil {
		logrus.Errorf("Failed to get all users: %v", err)
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	logrus.Infof("Retrieved %d users successfully", len(users))
	return users, nil
}
func (u userService) UpdatePhoto(ctx context.Context, id string, photoURL string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logrus.Infof("Updating photo for user ID: %s", id)

	err := u.userRepository.UpdatePhoto(ctx, id, photoURL)
	if err != nil {
		logrus.Errorf("Failed to update photo for user ID %s: %v", id, err)
		return fmt.Errorf("failed to update photo for user ID %s: %w", id, err)
	}

	logrus.Infof("Photo updated successfully for user ID: %s", id)
	return nil
}
