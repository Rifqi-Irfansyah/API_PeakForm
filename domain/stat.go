package domain

import (
	"api-peak-form/dto"
	"context"
)

type StatsService interface {
	GetStatsByUserID(ctx context.Context, userID string) (dto.StatsSummary, error)
}