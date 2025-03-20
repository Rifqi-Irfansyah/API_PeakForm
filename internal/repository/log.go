package repository

import "gorm.io/gorm"

type Log struct {
	db *gorm.DB
}

func NewLog(db *gorm.DB) *Log {
	return &Log{db: db}
}
