package domain

type ExerciseList struct {
	ScheduleID uint `gorm:"primaryKey"`
	ExerciseID uint `gorm:"primaryKey"`
	Set uint 		
	Repetition uint
}

func (ExerciseList) TableName() string {
    return "exercise_list"
}