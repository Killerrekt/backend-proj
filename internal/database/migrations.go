package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

func RunMigrations(db *gorm.DB) {
	log.Println("Running Migrations")

	err := db.AutoMigrate(&models.User{})
	err_invoice := db.AutoMigrate(&models.Invoice{})

	if err_invoice != nil {
		fmt.Println("Could not migrate Invoice")
		return
	}

	if err != nil {
		fmt.Println("Could not migrate")
		return
	}

	/*
		err := db.AutoMigrate(&models.User{}, &models.RefreshToken{})
		if err != nil {
			fmt.Println("Migration error")
			return
		}
	*/

	log.Println("🚀 Migrations completed")
}
