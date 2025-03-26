package dto

type CreateExerciseRequest struct {
	Name         string `json:"name" validate:"required"`
	Type         string `json:"type" validate:"required"`
	Muscle       string `json:"muscle" validate:"required"`
	Equipment    string `json:"equipment" validate:"required"`
	Difficulty   string `json:"difficulty" validate:"required"`
	Instructions string `json:"instructions" validate:"required"`
	Gif          string `json:"gif" validate:"required"`
}

type UpdateExerciseRequest struct {
	ID           uint   `json:"id" validate:"required"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Muscle       string `json:"muscle"`
	Equipment    string `json:"equipment"`
	Difficulty   string `json:"difficulty"`
	Instructions string `json:"instructions"`
	Gif          string `json:"gif"`
}
