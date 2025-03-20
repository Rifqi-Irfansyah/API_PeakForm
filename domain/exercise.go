package domain

type Exercise struct {
	ID           uint            `gorm:"primaryKey"`
	Name         string          `gorm:"type:varchar(100);not null"`
	Type         ExerciseType    `gorm:"type:ENUM('strength', 'cardio');not null"`
	Muscle       MuscleGroup     `gorm:"type:ENUM('abdominals', 'biceps', 'calves', 'chest', 'forearms', 'lats', 'lower_back', 'middle_back', 'neck', 'quadriceps', 'traps', 'triceps');not null"`
	Equipment    Equipment       `gorm:"type:ENUM('body_only', 'dumbbell');not null"`
	Difficulty   DifficultyLevel `gorm:"type:ENUM('beginner', 'intermediate', 'expert');not null"`
	Instructions string          `gorm:"type:TEXT;not null"`
	Gif          string          `gorm:"type:varchar(255);not null"`
}

type ExerciseRepository interface {
	// TYpe the method here
}
