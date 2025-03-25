package domain

import "context"

type User struct {
	ID        string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string     `gorm:"type:varchar(255);unique;not null"`
	Name      string     `gorm:"type:varchar(320);not null"`
	Password  string     `gorm:"type:char(60);not null"`
	Schedules []Schedule `gorm:"foreignKey:UID"`
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	Save(ctx context.Context, user User) error
	UpdatePassword(ctx context.Context, email string, password string) error
}
