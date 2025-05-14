package domain

import "context"

type User struct {
	ID        string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string     `gorm:"type:varchar(255);unique;not null"`
	Name      string     `gorm:"type:varchar(320);not null"`
	Password  string     `gorm:"type:char(60);not null"`
	Schedules []Schedule `gorm:"foreignKey:UID"`
	Point     int        `gorm:"default:0"`
	Streak    int        `gorm:"default:0"`
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id string) (User, error)
	Save(ctx context.Context, user User) error
	UpdatePassword(ctx context.Context, email string, password string) error
	UpdatePoint(ctx context.Context, id string, point int) error
	UpdateStreak(ctx context.Context, id string, streak int) error
}

type UserService interface {
	UpdatePoint(ctx context.Context, id string, difficulty DifficultyLevel, rep int, set int) (int, error)
	CheckStreak(ctx context.Context, id string) (int, error)
	UpdateStreak(ctx context.Context, id string) (int, error)
}
