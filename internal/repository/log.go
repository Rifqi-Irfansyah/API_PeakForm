package repository

import (
	"api-peak-form/domain"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type logRepository struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) domain.LogRepository {
	logrus.Info("Creating new LogRepository")
	return &logRepository{db: db}
}

func (l logRepository) Create(ctx context.Context, log domain.Log) error {
	logrus.Infof("Creating log for user ID: %s", log.UserID)
	if err := l.db.WithContext(ctx).Create(&log).Error; err != nil {
		logrus.Errorf("Failed to create log for user ID: %s, error: %v", log.UserID, err)
		return fmt.Errorf("failed to create log: %w", err)
	}
	logrus.Infof("Log created successfully for user ID: %s", log.UserID)
	return nil
}

func (l logRepository) FindByUserID(ctx context.Context, userID string) ([]domain.Log, error) {
	logrus.Infof("Fetching logs for user ID: %s", userID)
	var logs []domain.Log
	if err := l.db.WithContext(ctx).Preload("Exercise").Where("user_id = ?", userID).Find(&logs).Error; err != nil {
		logrus.Errorf("Failed to fetch logs for user ID %s: %v", userID, err)
		return nil, fmt.Errorf("failed to fetch logs for user ID %s: %w", userID, err)
	}
	logrus.Infof("Fetched %d logs for user ID: %s", len(logs), userID)
	return logs, nil
}

func (l logRepository) FindLastByUserID(ctx context.Context, userID string) (domain.Log, error) {
	var log domain.Log
	err := l.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("timestamp DESC").
		First(&log).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Infof("No logs found for user ID: %s", userID)
		return domain.Log{}, nil
	}

	if err != nil {
		logrus.Errorf("Failed to fetch last log for user ID %s: %v", userID, err)
		return domain.Log{}, fmt.Errorf("failed to fetch last log for user ID %s: %w", userID, err)
	}

	logrus.Infof("Fetched last log for user ID: %s", userID)
	return log, nil
}
