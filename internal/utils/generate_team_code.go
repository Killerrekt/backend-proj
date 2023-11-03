package utils

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

func GenerateUniqueTeamCode() (string, error) {
	var team models.Team
	for {
		code := fmt.Sprintf("%06s", uuid.New().String()[:6])
		result := database.DB.Where("code = ?", code).First(&team)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return code, nil
			}
			return "", result.Error
		}
	}
}
