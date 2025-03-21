package domain

import (
	"context"
	"time"
)

type Log struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null"`
	ScheduleID uint      `gorm:"not null"`
	Timestamp  time.Time `gorm:"autoCreateTime"`
	Status     string    `gorm:"type:varchar(50);not null"`
	User       User      `gorm:"foreignKey:UserID"`
	Schedule   Schedule  `gorm:"foreignKey:ScheduleID"`
}

type LogRepository interface {
	Create(ctx context.Context, log Log) error
	FindByUserID(ctx context.Context, userID uint) ([]Log, error)
}

type LogService interface {
	Create(ctx context.Context, log Log) error
	FindByUserID(ctx context.Context, userID uint) ([]Log, error)
}
