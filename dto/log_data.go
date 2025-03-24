package dto

import "time"

type LogRequest struct {
	UserID     string    `json:"user_id" validate:"required"`
	ExerciseID uint      `json:"exercise_id" validate:"required_with=numeric,gt=0"`
	Timestamp  time.Time `json:"timestamp" validate:"required"`
	Set        int       `json:"set" validate:"required_with=numeric,gt=0"`
	Repetition int       `json:"repetition" validate:"required,numeric,gt=0"`
}
