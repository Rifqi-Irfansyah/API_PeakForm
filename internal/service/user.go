package service

import (
	"api-peak-form/domain"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)


type userService struct {
	userRepository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) domain.UserService {
	return userService{userRepository: userRepository}
}

func (u userService) UpdatePoint(ctx context.Context, id string, point int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logrus.Infof("Total point for user ID: %s is %d", id, point)

	existingUser, err := u.userRepository.FindByID(ctx, id)
	if err != nil {
		logrus.Errorf("Failed to find user with ID %s: %v", id, err)
		return fmt.Errorf("failed to find user with ID %s: %w", id, err)
	}

	if existingUser.Streak == 0 {
		existingUser.Point += point
	} else {
		existingUser.Point += point * existingUser.Streak
	}


	logrus.Infof("User %s point updated to %d", id, existingUser.Point)

	err = u.userRepository.UpdatePoint(ctx, id, existingUser.Point)
	if err != nil {
		logrus.Errorf("Failed to update point and streak for user ID %s: %v", id, err)
		return fmt.Errorf("failed to update point and streak for user ID %s: %w", id, err)
	}

	logrus.Infof("Point and streak updated successfully for user ID: %s", id)
	return nil
}

func (u userService) UpdateStreak(ctx context.Context, id string, streak bool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logrus.Infof("Streak for user ID: %s is %d", id, streak)

	existingUser, err := u.userRepository.FindByID(ctx, id)
	if err != nil {
		logrus.Errorf("Failed to find user with ID %s: %v", id, err)
		return fmt.Errorf("failed to find user with ID %s: %w", id, err)
	}

	if streak == false {
		existingUser.Point = 0
	} else {
		existingUser.Point++
	}

	logrus.Infof("User %s streak updated to %d", id, existingUser.Streak)

	err = u.userRepository.UpdateStreak(ctx, id, existingUser.Streak)
	if err != nil {
		logrus.Errorf("Failed to update point and streak for user ID %s: %v", id, err)
		return fmt.Errorf("failed to update point and streak for user ID %s: %w", id, err)
	}

	logrus.Infof("Point and streak updated successfully for user ID: %s", id)
	return nil
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