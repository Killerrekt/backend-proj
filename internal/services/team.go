package services

import (
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

func FindTeamByID(id uint) (models.Team, error) {
	var team models.Team

	result := database.DB.Where("id = ?", id).First(&team)
	if result.Error != nil {
		return models.Team{}, result.Error
	}

	return team, nil
}
