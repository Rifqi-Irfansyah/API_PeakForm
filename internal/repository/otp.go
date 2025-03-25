package repository

import (
	"api-peak-form/domain"
	"context"
	"errors"
	"time"
)

type otpRepository struct {
	otpData map[string]struct {
		OTP       string
		ExpiresAt time.Time
	}
}

func NewOTPRepository() domain.OtpRepository {
	return &otpRepository{otpData: make(map[string]struct {
		OTP       string
		ExpiresAt time.Time
	})}
}

func (r *otpRepository) SaveOTP(ctx context.Context, email, otp string, expiresAt time.Time) error {
	r.otpData[email] = struct {
		OTP       string
		ExpiresAt time.Time
	}{OTP: otp, ExpiresAt: expiresAt}
	return nil
}

func (r *otpRepository) VerifyOTP(ctx context.Context, email, otp string) (bool, error) {
	data, exists := r.otpData[email]
	if !exists || data.ExpiresAt.Before(time.Now()) {
		return false, errors.New("OTP not found or expired")
	}
	return data.OTP == otp, nil
}

func (r *otpRepository) DeleteOTP(ctx context.Context, email string) error {
	delete(r.otpData, email)
	return nil
}
