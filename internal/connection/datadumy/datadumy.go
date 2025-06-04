package datadumy

import (
	"api-peak-form/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AddDefaultUser(db *gorm.DB) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	// Buat admin user
	admin := domain.User{
		ID:       "115dd593-1f58-454f-bd25-318cfd2b4819",
		Email:    "admin2@example.com",
		Name:     "Administrator",
		Password: string(hashedPassword),
		Point: 1000,
		Streak:  2,
	}

	user := domain.User{
		ID:       "115dd593-1f58-454f-bd25-318cfd2b4810",
		Email:    "user@example.com",
		Name:     "User",
		Password: string(hashedPassword),
		Point: 1500,
		Streak:  2,
	}

	user2 := domain.User{
		ID:       "115dd593-1f58-454f-bd25-318cfd2b4820",
		Email:    "fauzan@example.com",
		Name:     "Fauzan",
		Password: string(hashedPassword),
		Point: 2500,
		Streak:  4,
	}

	// Simpan ke database
	if err := db.Create(&admin).Error; err != nil {
		logrus.Errorf("Failed to create admin user: %v", err)
	} else {
		logrus.Info("Admin user created: admin@example.com")
	}

	if err := db.Create(&user).Error; err != nil {
		logrus.Errorf("Failed to create account user: %v", err)
	} else {
		logrus.Info("Account user created: user@example.com")
	}

	if err := db.Create(&user2).Error; err != nil {
		logrus.Errorf("Failed to create account user: %v", err)
	} else {
		logrus.Info("Account user created: fauzan@example.com")
	}
}