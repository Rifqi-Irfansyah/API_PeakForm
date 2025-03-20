package API_PeakForm

import (
	"api-peak-form/domain"
	"api-peak-form/internal/config"
	"api-peak-form/internal/connection"
	"log"
)

func main() {
	cnf := config.Get()
	dbConnection := connection.GetDatabase(cnf.Database)

	err := dbConnection.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully")

}
