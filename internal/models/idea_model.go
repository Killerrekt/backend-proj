package models

import "gorm.io/gorm"

type Idea struct {
	gorm.Model
	TeamID    uint   `json:"teamID" gorm:"unique"`
	Title     string `json:"title" gorm:"unique" validate:"required"`
	Desc      string `json:"desc" validate:"required"`
	Track     string `json:"tracK"`
	VideoLink string `json:"video_link"`
	FigmaLink string `json:"figma_link"`
	DriveLink string `json:"drive_link"`
}

type IdeaRequest struct {
	Title     string `json:"title,omitempty"`
	Desc      string `json:"desc,omitempty"`
	Track     string `json:"track,omitempty"`
	VideoLink string `json:"video_link,omitempty"`
	FigmaLink string `json:"figma_link,omitempty"`
	DriveLink string `json:"drive_link,omitempty"`
}
