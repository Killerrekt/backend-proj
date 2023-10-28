package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
<<<<<<< HEAD
	Name         string `json:"name"`
	Desc         string `json:"desc"`
	Githublink   string `json:"github"`
=======
	Name         string `json:"name" gorm:"unique" validate:"required"`
	Desc         string `json:"desc" validate:"required"`
	Githublink   string `json:"github" validate:"required"`
>>>>>>> 919f34b4797523e6af252a69b4a48588fd5be578
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
