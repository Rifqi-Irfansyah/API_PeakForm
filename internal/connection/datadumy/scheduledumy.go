package datadumy

import (
	"api-peak-form/domain"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddSchedules(db *gorm.DB) {
	// Hapus data di exercise_list terlebih dahulu karena memiliki foreign key ke schedules
	if err := db.Exec("DELETE FROM exercise_list").Error; err != nil {
		log.Println("Gagal menghapus data exercise_list:", err)
		return
	}
	log.Println("Semua data exercise_list telah dihapus")

	// Hapus semua data dalam tabel schedules setelah exercise_list dikosongkan
	if err := db.Exec("DELETE FROM schedules").Error; err != nil {
		log.Println("Gagal menghapus data schedule:", err)
		return
	}
	log.Println("Semua data schedule telah dihapus")


	// Ambil user pertama sebagai contoh``
	var user domain.User
	if err := db.First(&user).Error; err != nil {
		log.Println("User tidak ditemukan, pastikan ada user terlebih dahulu!")
		return
	}

	// Ambil beberapa exercise untuk digunakan dalam schedule
	var exercises []domain.Exercise
	if err := db.Limit(3).Find(&exercises).Error; err != nil {
		log.Println("Exercise tidak ditemukan, pastikan ada exercise terlebih dahulu!")
		return
	}

	var exercisesMonday, exercisesWednesday, exercisesFriday []domain.Exercise
	db.Limit(3).Find(&exercisesMonday)
	db.Offset(3).Limit(3).Find(&exercisesWednesday)
	db.Offset(6).Limit(3).Find(&exercisesFriday)

	// Data dummy untuk Schedule
	schedules := []domain.Schedule{
		{
			ID:        uint(uuid.New().ID()),
			UserID:    "7047aa3b-67fe-458f-908f-0053f137f491",
			Day:       1, // Senin
			Exercises: exercisesMonday,
		},
		{
			ID:        uint(uuid.New().ID()),
			UserID:    "7047aa3b-67fe-458f-908f-0053f137f491",
			Day:       3, // Rabu
			Exercises: exercisesWednesday,
		},
		{
			ID:        uint(uuid.New().ID()),
			UserID:    "7047aa3b-67fe-458f-908f-0053f137f491",
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