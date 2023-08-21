package database

import (
	"gorm.io/gorm"
	"log"
)

func RunMigrations(db *gorm.DB) {
	log.Println("Running Migrations")

	/*
		err := db.AutoMigrate(&models.User{}, &models.RefreshToken{})
		if err != nil {
			fmt.Println("Migration error")
			return
		}
	*/

	log.Println("🚀 Migrations completed")
}
