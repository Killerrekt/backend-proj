package models

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Users   []User  `gorm:"foreignKey:Email"`
	Project Project `gorm:"foreignKey:Name"`
}
