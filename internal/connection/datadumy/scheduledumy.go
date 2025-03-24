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
			ID:        uint(uuid.New().ID()),
			Type:		"strength",
			Day:       1, // Senin
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

func AddUserSchedules(db *gorm.DB){
	var user domain.User
	if err := db.First(&user).Error; err != nil {
		log.Println("User tidak ditemukan, pastikan ada user terlebih dahulu!")
		return
	}

	var schedule domain.Schedule
	if err := db.First(&schedule).Error; err != nil {
		log.Println("Schedule tidak ditemukan, pastikan ada schedule terlebih dahulu!")
		return
	}


	type UserSchedule struct {
		UserID     string `gorm:"column:user_id"`
		ScheduleID uint   `gorm:"column:schedule_id"`
	}
	
	userschedule := UserSchedule{
		UserID: user.ID,
		ScheduleID: schedule.ID,
	}

	log.Println("ID Schedule: ", schedule.ID)
	if err := db.Table("user_schedules").Create(&userschedule).Error; err != nil {
		log.Println("Gagal menambahkan data user schedule:", err)
	} else {
		log.Println("Schedule berhasil ditambahkan untuk user:", user.ID)
	}
	log.Println("Schedule ditambahkan untuk user:", user.ID)
}