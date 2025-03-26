package repository

import (
	"context"
	"api-peak-form/domain"

	"gorm.io/gorm"
)

type exerciseRepository struct {
	db *gorm.DB
}

func NewExercise(db *gorm.DB) domain.ExerciseRepository {
	return &exerciseRepository{db}
}

func (r *exerciseRepository) Create(ctx context.Context, exercise *domain.Exercise) error {
	return r.db.WithContext(ctx).Create(exercise).Error
}

func (r *exerciseRepository) GetAll(ctx context.Context) ([]domain.Exercise, error) {
	var exercises []domain.Exercise
	err := r.db.WithContext(ctx).Find(&exercises).Error
	return exercises, err
}

func (r *exerciseRepository) GetByID(ctx context.Context, id uint) (*domain.Exercise, error) {
	var exercise domain.Exercise
	err := r.db.WithContext(ctx).First(&exercise, id).Error
	if err != nil {
		return nil, err
	}
	return &exercise, nil
}

func (r *exerciseRepository) Update(ctx context.Context, exercise *domain.Exercise) error {
	return r.db.WithContext(ctx).Save(exercise).Error
}

func (r *exerciseRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Exercise{}, id).Error
}
