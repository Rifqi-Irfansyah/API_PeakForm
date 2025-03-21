package dto

import "time"

type LogRequest struct {
	UserID    uint      `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}
