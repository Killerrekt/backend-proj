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
	err_project := db.AutoMigrate(&models.Project{})
	err_team := db.AutoMigrate(&models.Team{})
	err_enquiry := db.AutoMigrate(&models.Enquiry{})
	err_idea := db.AutoMigrate(&models.Idea{})

	if err_invoice != nil {
		fmt.Println("Could not migrate Invoice")
		return
	}

	if err != nil {
		fmt.Println("Could not migrate")
		return
	}

	if err_project != nil {
		fmt.Println("could not migrate projects")
		return
	}

	if err_team != nil {
		fmt.Println("Could not migrate teams")
		return
	}

	if err_enquiry != nil {
		fmt.Println("Could not migrate enquiries")
		return
	}

	if err_idea != nil {
		fmt.Println("Could not migrate idea")
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
