package domain

type ExerciseType string

const (
	Strength ExerciseType = "strength"
	Cardio   ExerciseType = "cardio"
)

type MuscleGroup string

const (
	Abdominals MuscleGroup = "abdominals"
	Biceps     MuscleGroup = "biceps"
	Calves     MuscleGroup = "calves"
	Chest      MuscleGroup = "chest"
	Forearms   MuscleGroup = "forearms"
	Lats       MuscleGroup = "lats"
	LowerBack  MuscleGroup = "lower_back"
	MiddleBack MuscleGroup = "middle_back"
	Neck       MuscleGroup = "neck"
	Quadriceps MuscleGroup = "quadriceps"
	Traps      MuscleGroup = "traps"
	Triceps    MuscleGroup = "triceps"
)

type Equipment string

const (
	BodyOnly Equipment = "body_only"
	Dumbbell Equipment = "dumbbell"
)

type DifficultyLevel string

const (
	Beginner     DifficultyLevel = "beginner"
	Intermediate DifficultyLevel = "intermediate"
	Expert       DifficultyLevel = "expert"
)
