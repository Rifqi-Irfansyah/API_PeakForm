package datadumy

import (
	"api-peak-form/domain"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddSchedules(db *gorm.DB) {
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
			UID:		"115dd593-1f58-454f-bd25-318cfd2b4819",
			ID:        	uint(uuid.New().ID()),
			Type:		"strength",
			Day:       	1,
			Exercises: exercisesMonday,
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