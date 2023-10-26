package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Name         string `json:"name" gorm:"unique" validate:"required"`
	Desc         string `json:"desc" validate:"required"`
	Githublink   string `json:"github" validate:"required"`
	FigmaLink    string `json:"figma"`
	VideoLink    string `json:"video"`
	DriveLink    string `json:"drive"`
	ProjectTrack string `json:"project_track" validate:"required"`
	IsFinal      bool   `json:"isfinal"`
	TeamID       uint   `json:"teamID"`
}

type CreateProject struct {
	//Name         string `json:"name"`
	Desc       string `json:"desc"`
	Githublink string `json:"github"`
	FigmaLink  string `json:"figma"`
	VideoLink  string `json:"video"`
	DriveLink  string `json:"drive"`
	//ProjectTrack string `json:"project_track"`
}

type GetProject struct {
	UserID  uint `json:"userID"`
	TeamID  uint `json:"teamID"`
	Isfinal bool `json:"isfinal"`
}
