package domain

import (
	"api-peak-form/dto"
	"context"
)

type AuthService interface {
	Login(ctx context.Context, data dto.AuthRequest) (dto.AuthResponse, error)
}
