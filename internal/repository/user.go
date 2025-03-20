package repository

import (
	"api-peak-form/domain"
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (u userRepository) FindByEmail(ctx context.Context, email string) (user domain.User, err error) {
	err = u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("User tidak ditemukan untuk email:", email)
			return user, nil
		}
		log.Println("Error query user:", err)
		return user, err
	}
	log.Println("User ditemukan:", user.Email)
	return user, nil
}
