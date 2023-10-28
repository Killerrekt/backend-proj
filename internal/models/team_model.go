package models

import "gorm.io/gorm"

type Team struct {
	gorm.Model
	Name string
	// TeamCode string `gorm:"uniqueIndex"`
	// Members      []User  `gorm:"one2many:team_members;"`
	ProjectID    uint
	Round        int
	IdeaID       uint
	LeaderID     uint
	MembersCount int     //`gorm:"-"` // not saved in DB
	TeamID       uint    `gorm:"primaryKey;unique"`
	Users        []User  `gorm:"foreignKey:TeamID;references:TeamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Project      Project `gorm:"foreignKey:TeamID;references:TeamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
