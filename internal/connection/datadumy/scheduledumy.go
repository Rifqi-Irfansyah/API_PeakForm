package datadumy

import (
	"api-peak-form/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func AddSchedules(db *gorm.DB) {
	logrus.Info("Starting to add schedules")

	var exercises []domain.Exercise
	if err := db.Limit(3).Find(&exercises).Error; err != nil {
		logrus.Warn("Exercise tidak ditemukan, pastikan ada exercise terlebih dahulu!")
		return
	}

	var exercisesMonday, exercisesWednesday, exercisesFriday []domain.Exercise
	db.Limit(3).Find(&exercisesMonday)
	db.Offset(3).Limit(3).Find(&exercisesWednesday)
	db.Offset(6).Limit(3).Find(&exercisesFriday)

	schedules := []domain.Schedule{
		{
			UID:       "115dd593-1f58-454f-bd25-318cfd2b4819",
			ID:        uint(uuid.New().ID()),
			Type:      "strength",
			Day:       1,
			Exercises: exercisesMonday,
		},
	}

	// Simpan ke database
	for _, schedule := range schedules {
		if err := db.Create(&schedule).Error; err != nil {
			logrus.Errorf("Gagal menambahkan data schedule: %v", err)
		} else {
			logrus.Infof("Schedule ditambahkan untuk hari: %d", schedule.Day)
		}
	}

	logrus.Info("Finished adding schedules")
}
