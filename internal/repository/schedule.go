package repository

import "gorm.io/gorm"

type Schedule struct {
	db *gorm.DB
}

func NewSchedule(db *gorm.DB) *Schedule {
	return &Schedule{db: db}
}
