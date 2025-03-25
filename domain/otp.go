package domain

import (
	"context"
	"time"
)

type OtpRepository interface {
	SaveOTP(ctx context.Context, email string, otp string, expiresAt time.Time) error
	VerifyOTP(ctx context.Context, email string, otp string) (bool, error)
	DeleteOTP(ctx context.Context, email string) error
}
