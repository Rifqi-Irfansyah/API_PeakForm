package repository

import (
	"api-peak-form/domain"
	"context"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type exerciseRepository struct {
	db *gorm.DB
}

func NewExercise(db *gorm.DB) domain.ExerciseRepository {
	logrus.Info("Creating new ExerciseRepository")
	return &exerciseRepository{db}
}

func (r *exerciseRepository) Create(ctx context.Context, exercise *domain.Exercise) error {
	err := r.db.WithContext(ctx).Create(exercise).Error
	if err != nil {
		logrus.Errorf("Error creating exercise: %v", err)
		return err
	}
	logrus.Infof("Exercise created successfully: %v", exercise)
	return nil
}

func (r *exerciseRepository) GetAll(ctx context.Context) ([]domain.Exercise, error) {
	var exercises []domain.Exercise
	err := r.db.WithContext(ctx).Find(&exercises).Error
	if err != nil {
		logrus.Errorf("Error fetching exercises: %v", err)
		return nil, err
	}
	logrus.Infof("Fetched %d exercises", len(exercises))
	return exercises, nil
}

func (r *exerciseRepository) GetByID(ctx context.Context, id uint) (*domain.Exercise, error) {
	var exercise domain.Exercise
	err := r.db.WithContext(ctx).First(&exercise, id).Error
	if err != nil {
		logrus.Errorf("Error fetching exercise with ID %d: %v", id, err)
		return nil, err
	}
	logrus.Infof("Exercise fetched successfully with ID %d: %v", id, exercise)
	return &exercise, nil
}

func (r *exerciseRepository) Update(ctx context.Context, exercise *domain.Exercise) error {
	err := r.db.WithContext(ctx).Save(exercise).Error
	if err != nil {
		logrus.Errorf("Error updating exercise: %v", err)
		return err
	}
	logrus.Infof("Exercise updated successfully: %v", exercise)
	return nil
}

func (r *exerciseRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&domain.Exercise{}, id).Error
	if err != nil {
		logrus.Errorf("Error deleting exercise with ID %d: %v", id, err)
		return err
	}
	logrus.Infof("Exercise deleted successfully with ID %d", id)
	return nil
}
