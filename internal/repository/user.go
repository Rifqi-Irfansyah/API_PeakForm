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
			log.Println("User not found for email:", email)
			return user, nil
		}
		log.Println("Error querying user:", err)
		return user, err
	}
	log.Println("User found:", user.Email)
	return user, nil
}

func (u userRepository) Save(ctx context.Context, user domain.User) error {
	err := u.db.WithContext(ctx).Save(&user).Error
	if err != nil {
		log.Println("Error while saving user:", err)
		return err
	}
	log.Println("User successfully saved:", user.Email)
	return nil
}
