package domain

import "time"

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
	Create(log Log) error
	FindAll() ([]Log, error)
}
