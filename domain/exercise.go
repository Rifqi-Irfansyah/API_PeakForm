package domain

type Exercise struct {
	ID           uint            `gorm:"primaryKey"`
	Name         string          `gorm:"type:varchar(100);not null"`
	Type         ExerciseType    `gorm:"type:exercise_type;not null"`
	Muscle       MuscleGroup     `gorm:"type:muscle_group;not null"`
	Equipment    Equipment       `gorm:"type:equipment;not null"`
	Difficulty   DifficultyLevel `gorm:"type:difficulty_level;not null"`
	Instructions string          `gorm:"type:TEXT;not null"`
	Gif          string          `gorm:"type:varchar(255);not null"`
}

type ExerciseRepository interface {
	// TYpe the method here
}
