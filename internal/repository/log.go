package repository

import (
	"api-peak-form/domain"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type logRepository struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) domain.LogRepository {
	return &logRepository{db: db}
}

func (l logRepository) Create(ctx context.Context, log domain.Log) error {
	if err := l.db.WithContext(ctx).Create(&log).Error; err != nil {
		return fmt.Errorf("failed to create log: %w", err)
	}
	return nil
}

func (l logRepository) FindByUserID(ctx context.Context, userID uint) ([]domain.Log, error) {
	var logs []domain.Log
	if err := l.db.WithContext(ctx).Where("user_id = ?", userID).Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch logs for user ID %d: %w", userID, err)
	}
	return logs, nil
}
