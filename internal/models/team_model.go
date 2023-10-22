package models

type Team struct {
	TeamID  uint    `gorm:"primaryKey;unique"`
	Users   []User  `gorm:"foreignKey:TeamId;references:TeamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Project Project `gorm:"foreignKey:TeamID;references:TeamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
