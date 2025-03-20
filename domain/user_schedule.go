package domain

type UserSchedule struct {
	UserID     uint `gorm:"primaryKey"`
	ScheduleID uint `gorm:"primaryKey"`
}
