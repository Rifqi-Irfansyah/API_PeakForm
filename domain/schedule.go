package domain

type Schedule struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	User      User       `gorm:"foreignKey:UserID"`
	Day       int        `gorm:"not null; check:day >= 1 AND day <= 7"`
	Exercises []Exercise `gorm:"many2many:exercise_list"`
}

type ScheduleRepository interface {
	// Type method here
}
