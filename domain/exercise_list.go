package domain

type ExerciseList struct {
	ScheduleID uint `gorm:"primaryKey"`
	ExerciseID uint `gorm:"primaryKey"`
	Set uint
	Repetition uint

	Schedule Schedule `gorm:"foreignKey:ScheduleID;references:ID;constraint:OnDelete:CASCADE"`
	Exercise Exercise `gorm:"foreignKey:ExerciseID;references:ID;constraint:OnDelete:CASCADE"`
}

func (ExerciseList) TableName() string {
    return "exercise_list"
}