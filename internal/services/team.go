package services

import (
	"fmt"

	"gorm.io/gorm/clause"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/utils"
)

func FindTeamByID(id uint) (models.Team, error) {
	var team models.Team

	result := database.DB.Where("team_id = ?", id).Preload(clause.Associations).First(&team)
	if result.Error != nil {
		return models.Team{}, result.Error
	}

	return team, nil
}

func FindTeamByName(name string) (models.Team, error) {
	var team models.Team
	if err := database.DB.Where("name = ?", name).
		Preload("Users").
		Preload("Project").
		Preload("Idea").
		First(&team).Error; err != nil {
		return models.Team{}, err
	}
	return team, nil
}

func FindTeamByCode(code string) (models.Team, error) {
	var team models.Team
	result := database.DB.Where("code = ?", code).
		Preload("Users").
		Preload("Project").
		Preload("Idea").
		First(&team)
	if result.Error != nil {
		return models.Team{}, result.Error
	}

	return team, nil
}

func DeleteTeamByID(id uint) error {
	result, err := FindTeamByID(id)
	if err != nil {
		return err
	}

	for _, user := range result.Users {
		user.TeamID = 0
		if err := database.DB.Save(&user).Error; err != nil {
			return err
		}

		err := utils.SendMail(
			"Team Deleted",
			fmt.Sprintf("The team %s has been deleted", result.Name),
			user.Email,
		)
		if err != nil {
			return err
		}
	}

	if err := database.DB.Unscoped().Delete(&result).Error; err != nil {
		return err
	}

	return nil
}
