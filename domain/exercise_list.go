package domain

type ExerciseList struct {
	ScheduleID uint `gorm:"primaryKey"`
	ExerciseID uint `gorm:"primaryKey"`
}
