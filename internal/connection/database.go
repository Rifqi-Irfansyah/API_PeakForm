package connection

import (
	"api-peak-form/internal/config"
	"api-peak-form/internal/connection/migration"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func GetDatabase(conf config.Database) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		conf.Host,
		conf.Port,
		conf.User,
		conf.Pass,
		conf.Name,
		conf.Tz,
	)

	


	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get generic database connection: ", err.Error())
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(0)

	if err := migration.Migrate(db); err != nil {
		log.Fatal("Migration error:", err)
	}

	log.Println("Database connection established!")

	return db
}
