package domain

import (
	"api-peak-form/dto"
	"context"
)

type AuthService interface {
	Login(ctx context.Context, data dto.AuthRequest) (dto.AuthResponse, error)
	Register(ctx context.Context, data dto.RegisterRequest) error
	VerifyRegisterOTP(ctx context.Context, data dto.VerifyOTPRequest) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, data dto.ResetPasswordRequest) error
	ChangePassword(ctx context.Context, data dto.ChangePasswordRequest) error
	CheckToken(ctx context.Context, token string) error
	GetUserByToken(ctx context.Context, token string) (dto.AuthResponse, error)
}
