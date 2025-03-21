package domain

type ExerciseList struct {
	ScheduleID uint `gorm:"primaryKey"`
	ExerciseID uint `gorm:"primaryKey"`
	Set        int  `gorm:"not null"`
	Repetition int  `gorm:"not null"`
}
