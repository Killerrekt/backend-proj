package models

import "gorm.io/gorm"

type Idea struct {
	gorm.Model
	TeamID    uint   `json:"teamID"`
	Title     string `json:"title" gorm:"unique" validate:"required"`
	Desc      string `json:"desc" validate:"required"`
	VideoLink string `json:"video_link"`
	FigmaLink string `json:"figma_link"`
	DriveLink string `json:"drive_link"`
}
