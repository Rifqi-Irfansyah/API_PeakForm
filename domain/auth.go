package domain

import (
	"api-peak-form/dto"
	"context"
)

type Auth interface {
	Login(ctx context.Context, data dto.AuthRequest) (dto.AuthResponse, error)
}
