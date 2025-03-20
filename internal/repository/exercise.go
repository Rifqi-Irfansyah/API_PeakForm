package repository

import "gorm.io/gorm"

type Exercise struct {
	db *gorm.DB
}

func NewExercise(db *gorm.DB) *Exercise {
	return &Exercise{db: db}
}
