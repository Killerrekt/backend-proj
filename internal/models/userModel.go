package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"         gorm:"unique"`
	Role        string    `json:"role"          gorm:"default:user"`
	Password    string    `json:"password"`
	Gender      string    `json:"gender"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Bio         string    `json:"bio"`
	TeamId      int       `json:"team_id"` // TODO: Link to team model
	IsLeader    bool      `json:"is_leader"`
	IsApproved  bool      `json:"is_approved"`
	PhoneNumber string    `json:"phone_number"`
	College     string    `json:"college"`
	Github      string    `json:"github"`
}

type CreateUser struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"` // considering "YYYY/MM/DD format"
	Bio         string `json:"bio"`
	PhoneNumber string `json:"phone_number"`
	College     string `json:"college"`
	Github      string `json:"github"`
}
