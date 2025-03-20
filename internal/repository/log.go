package repository

import (
	"api-peak-form/domain"
	"context"
	"gorm.io/gorm"
)

type logRepository struct {
	db *gorm.DB
}

func NewLog(db *gorm.DB) domain.LogRepository {
	return &logRepository{db: db}
}

func (l logRepository) Create(ctx context.Context, log domain.Log) error {
	return l.db.WithContext(ctx).Create(&log).Error
}

func (l logRepository) FindAll(ctx context.Context) ([]domain.Log, error) {
	var logs []domain.Log
	return logs, l.db.WithContext(ctx).Find(&logs).Error
}
