package dto

import "time"

type LogRequest struct {
	UserID     string    `json:"user_id" validate:"required;numeric,gt=0"`
	ExerciseID uint      `json:"exercise_id" validate:"required;numeric,gt=0"`
	Timestamp  time.Time `json:"timestamp" validate:"required;date"`
	Set        int       `json:"set" validate:"required;numeric,gt=0"`
	Repetition int       `json:"repetition" validate:"required;numeric,gt=0"`
}
