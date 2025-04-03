package datadumy

import (
	"api-peak-form/domain"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func AddExercise(db *gorm.DB) {
	// Data dummy untuk tabel Exercise
	exercises := []domain.Exercise{
		{
			ID: 100,
			Name:         "Bicep Curl",
			Type:         "strength",
			Muscle:       "biceps",
			Equipment:    "dumbbell",
			Difficulty:   "intermediate",
			Instructions: "Pegang dumbbell di kedua tangan, angkat ke arah bahu sambil mempertahankan siku tetap di sisi tubuh, lalu turunkan perlahan.",
			Image:          "https://example.com/bicepcurl.gif",
		},
		{
			ID: 101,
			Name:         "Deadlift",
			Type:         "strength",
			Muscle:       "abdominals",
			Equipment:    "dumbbell",
			Difficulty:   "intermediate",
			Instructions: "Berdiri dengan kaki selebar bahu, pegang barbell dengan genggaman overhand, angkat ke atas dengan menggunakan pinggul dan punggung bawah.",
			Image:          "https://example.com/deadlift.gif",
		},
	}

	for _, exercise := range exercises {
		if err := db.Create(&exercise).Error; err != nil {
			logrus.Errorf("Failed to add exercise data: %v", err)
		} else {
			logrus.Infof("Exercise added: %s", exercise.Name)
		}
	}
}
