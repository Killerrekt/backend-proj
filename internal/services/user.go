package services

import (
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := database.DB.Find(&users)
	if result.Error != nil {
		return []models.User{}, result.Error
	}

	return users, nil
}

func FindUserByEmail(email string) (models.User, error) {
	var user models.User

	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	return user, nil
}

func FindUserByID(id uint) (models.User, error) {
	var user models.User

	result := database.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	return user, nil
}
