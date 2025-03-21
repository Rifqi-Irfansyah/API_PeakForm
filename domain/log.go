package domain

import (
	"context"
	"time"
)

type Log struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null"`
	ExerciseID uint      `gorm:"not null"`
	Timestamp  time.Time `gorm:"not null"`
	User       User      `gorm:"foreignKey:UserID"`
	Exercise   Exercise  `gorm:"foreignKey:ExerciseID"`
	Set        int       `gorm:"not null"`
	Repetition int       `gorm:"not null"`
}

type LogRepository interface {
	Create(ctx context.Context, log Log) error
	FindAll(ctx context.Context) ([]Log, error)
}
