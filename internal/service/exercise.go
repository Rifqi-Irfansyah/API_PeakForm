package service

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"context"
	"github.com/sirupsen/logrus"
)

type exerciseService struct {
	repo domain.ExerciseRepository
}

func NewExerciseService(repo domain.ExerciseRepository) domain.ExerciseService {
	return &exerciseService{repo}
}

func (s *exerciseService) CreateExercise(ctx context.Context, req dto.CreateExerciseRequest) error {
	logrus.Infof("Creating exercise with name: %s", req.Name)
	exercise := &domain.Exercise{
		Name:         req.Name,
		Type:         domain.ExerciseType(req.Type),
		Muscle:       domain.MuscleGroup(req.Muscle),
		Equipment:    domain.Equipment(req.Equipment),
		Difficulty:   domain.DifficultyLevel(req.Difficulty),
		Instructions: req.Instructions,
		Image:        req.Image,
	}
	err := s.repo.Create(ctx, exercise)
	if err != nil {
		logrus.Errorf("Failed to create exercise: %v", err)
		return err
	}
	logrus.Infof("Exercise created successfully with name: %s", req.Name)
	return nil
}

func (s *exerciseService) GetExercises(ctx context.Context) ([]domain.Exercise, error) {
	logrus.Info("Fetching all exercises")
	exercises, err := s.repo.GetAll(ctx)
	if err != nil {
		logrus.Errorf("Failed to fetch exercises: %v", err)
		return nil, err
	}
	logrus.Infof("Fetched %d exercises", len(exercises))
	return exercises, nil
}

func (s *exerciseService) GetExerciseByID(ctx context.Context, id uint) (domain.Exercise, error) {
	logrus.Infof("Fetching exercise with ID: %d", id)
	exercise, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logrus.Errorf("Failed to fetch exercise with ID %d: %v", id, err)
		return domain.Exercise{}, err
	}
	logrus.Infof("Fetched exercise with ID: %d", id)
	return *exercise, nil
}

func (s *exerciseService) UpdateExercise(ctx context.Context, req dto.UpdateExerciseRequest) error {
	logrus.Infof("Updating exercise with ID: %d", req.ID)

	existingExercise, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		logrus.Errorf("Failed to fetch exercise with ID %d: %v", req.ID, err)
		return err
	}

	if req.Name != "" {
		existingExercise.Name = req.Name
	}
	if req.Type != "" {
		existingExercise.Type = domain.ExerciseType(req.Type)
	}
	if req.Muscle != "" {
		existingExercise.Muscle = domain.MuscleGroup(req.Muscle)
	}
	if req.Equipment != "" {
		existingExercise.Equipment = domain.Equipment(req.Equipment)
	}
	if req.Difficulty != "" {
		existingExercise.Difficulty = domain.DifficultyLevel(req.Difficulty)
	}
	if req.Instructions != "" {
		existingExercise.Instructions = req.Instructions
	}
	if req.Image != "" {
		existingExercise.Image = req.Image
	}

	err = s.repo.Update(ctx, existingExercise)
	if err != nil {
		logrus.Errorf("Failed to update exercise with ID %d: %v", req.ID, err)
		return err
	}

	logrus.Infof("Exercise updated successfully with ID: %d", req.ID)
	return nil
}

func (s *exerciseService) DeleteExercise(ctx context.Context, id uint) error {
	logrus.Infof("Attempting to delete exercise with ID: %d", id)
	err := s.repo.Delete(ctx, id)
	if err != nil {
		logrus.Errorf("Failed to delete exercise with ID %d: %v", id, err)
		return err
	}
	logrus.Infof("Exercise deleted successfully with ID: %d", id)
	return nil
}
