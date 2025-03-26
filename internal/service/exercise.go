package service

import (
	"context"
	"api-peak-form/domain"
	"api-peak-form/dto"
)

type exerciseService struct {
	repo domain.ExerciseRepository
}

func NewExerciseService(repo domain.ExerciseRepository) domain.ExerciseService {
	return &exerciseService{repo}
}

func (s *exerciseService) CreateExercise(ctx context.Context, req dto.CreateExerciseRequest) error {
	exercise := &domain.Exercise{
		Name:         req.Name,
		Type:         domain.ExerciseType(req.Type),
		Muscle:       domain.MuscleGroup(req.Muscle),
		Equipment:    domain.Equipment(req.Equipment),
		Difficulty:   domain.DifficultyLevel(req.Difficulty),
		Instructions: req.Instructions,
		Gif:          req.Gif,
	}
	return s.repo.Create(ctx, exercise)
}

func (s *exerciseService) GetExercises(ctx context.Context) ([]domain.Exercise, error) {
	return s.repo.GetAll(ctx)
}

func (s *exerciseService) GetExerciseByID(ctx context.Context, id uint) (domain.Exercise, error) {
	exercise, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Exercise{}, err
	}
	return *exercise, nil 
}

func (s *exerciseService) UpdateExercise(ctx context.Context, req dto.UpdateExerciseRequest) error {
	
	existingExercise, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
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
	if req.Gif != "" {
		existingExercise.Gif = req.Gif
	}

	return s.repo.Update(ctx, existingExercise)
}


func (s *exerciseService) DeleteExercise(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}