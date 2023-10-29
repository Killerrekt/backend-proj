package models

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	TeamID   uint   `gorm:"primaryKey;unique"`
	Code     string `gorm:"unique"                                                                           json:"code"`
	Name     string
	Round    int
	LeaderID uint
	Users    []User  `gorm:"foreignKey:TeamID;references:TeamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Project  Project `gorm:"foreignKey:TeamID;references:TeamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Idea     Idea    `gorm:"foreignKey:TeamID;references:TeamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
