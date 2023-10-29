package utils

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/services"
)

func GenerateUniqueTeamCode() (string, error) {
	for {
		code := fmt.Sprintf("%06s", uuid.New().String()[:6])
		_, err := services.FindTeamByCode(code)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return code, nil
			}
			return "", err
		}
	}
}
