package datadumy

import (
	"api-peak-form/domain"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddSchedules(db *gorm.DB) {
	var user domain.User
	if err := db.First(&user).Error; err != nil {
		log.Println("User tidak ditemukan, pastikan ada user terlebih dahulu!")
		return
	}

	var exercises []domain.Exercise
	if err := db.Limit(3).Find(&exercises).Error; err != nil {
		log.Println("Exercise tidak ditemukan, pastikan ada exercise terlebih dahulu!")
		return
	}

	var exercisesMonday, exercisesWednesday, exercisesFriday []domain.Exercise
	db.Limit(3).Find(&exercisesMonday)
	db.Offset(3).Limit(3).Find(&exercisesWednesday)
	db.Offset(6).Limit(3).Find(&exercisesFriday)

	schedules := []domain.Schedule{
		{
			ID:        uint(uuid.New().ID()),
			UserID:    user.ID,
			Day:       1, // Senin
			Exercises: exercisesMonday,
		},
		{
			ID:        uint(uuid.New().ID()),
			UserID:    user.ID,
			Day:       3, // Rabu
			Exercises: exercisesWednesday,
		},
		{
			ID:        uint(uuid.New().ID()),
			UserID:    user.ID,
			Day:       5, // Jumat
			Exercises: exercisesFriday,
		},
	}

	// Simpan ke database
	for _, schedule := range schedules {
		if err := db.Create(&schedule).Error; err != nil {
			log.Println("Gagal menambahkan data schedule:", err)
		} else {
			log.Println("Schedule ditambahkan untuk hari:", schedule.Day)
		}
	}
}