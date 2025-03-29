package repository

import (
	"api-peak-form/domain"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"time"
)

type otpRepository struct {
	otpData map[string]struct {
		OTP       string
		ExpiresAt time.Time
	}
}

func NewOTPRepository() domain.OtpRepository {
	logrus.Info("Creating new OTPRepository")
	return &otpRepository{otpData: make(map[string]struct {
		OTP       string
		ExpiresAt time.Time
	})}
}

func (r *otpRepository) SaveOTP(ctx context.Context, email, otp string, expiresAt time.Time) error {
	logrus.Infof("Saving OTP for email: %s", email)
	r.otpData[email] = struct {
		OTP       string
		ExpiresAt time.Time
	}{OTP: otp, ExpiresAt: expiresAt}
	logrus.Infof("OTP saved successfully for email: %s", email)
	return nil
}

func (r *otpRepository) VerifyOTP(ctx context.Context, email, otp string) (bool, error) {
	logrus.Infof("Verifying OTP for email: %s", email)
	data, exists := r.otpData[email]
	if !exists || data.ExpiresAt.Before(time.Now()) {
		logrus.Warnf("OTP not found or expired for email: %s", email)
		return false, errors.New("OTP not found or expired")
	}
	isValid := data.OTP == otp
	if isValid {
		logrus.Infof("OTP verified successfully for email: %s", email)
	} else {
		logrus.Warnf("Invalid OTP for email: %s", email)
	}
	return isValid, nil
}

func (r *otpRepository) DeleteOTP(ctx context.Context, email string) error {
	logrus.Infof("Deleting OTP for email: %s", email)
	delete(r.otpData, email)
	logrus.Infof("OTP deleted successfully for email: %s", email)
	return nil
}
