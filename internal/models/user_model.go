package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"         gorm:"unique"`
	Role         string    `json:"role"          gorm:"default:user"`
	Password     string    `json:"password"`
	Gender       string    `json:"gender"`
	Country      string    `json:"country"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	Bio          string    `json:"bio"`
	TeamID       uint      `json:"team_id"`
	IsLeader     bool      `json:"is_leader"     gorm:"default:false"`
	IsApproved   bool      `json:"is_approved"   gorm:"default:false"`
	IsVerified   bool      `json:"is_verified"   gorm:"default:false"`
	IsBanned     bool      `json:"is_banned"     gorm:"default:false"`
	IsPaid       bool      `json:"is_paid"       gorm:"default:false"`
	PhoneNumber  string    `json:"phone_number"`
	College      string    `json:"college"`
	Github       string    `json:"github"`
	TokenVersion int       `json:"token_version" gorm:"default:0"`
	Invoice      []Invoice `                     gorm:"foreignKey:UserID;References:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type CreateUser struct {
	FirstName   string `json:"first_name"    validate:"required"`
	LastName    string `json:"last_name"     validate:"required"`
	Email       string `json:"email"         validate:"required"`
	Password    string `json:"password"      validate:"required"`
	Gender      string `json:"gender"        validate:"required"`
	DateOfBirth string `json:"date_of_birth" validate:"required"` // considering "YYYY/MM/DD format"
	Bio         string `json:"bio"           validate:"required"`
	PhoneNumber string `json:"phone_number"  validate:"required"`
	College     string `json:"college"       validate:"required"`
	Github      string `json:"github"        validate:"required"`
	Country     string `json:"country"       validate:"required"`
}

type UserProfile struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
	Bio         string `json:"bio"`
	PhoneNumber string `json:"phone_number"`
	College     string `json:"college"`
	Github      string `json:"github"`
	Country     string `json:"country"`
	Team        Team   `json:"team"`
}

type UpdateUser struct {
	FirstName   string `json:"first_name"    validate:"required"`
	LastName    string `json:"last_name"     validate:"required"`
	Gender      string `json:"gender"        validate:"required"`
	DateOfBirth string `json:"date_of_birth" validate:"required"` // considering "YYYY/MM/DD format"
	Bio         string `json:"bio"           validate:"required"`
	PhoneNumber string `json:"phone_number"  validate:"required"`
	College     string `json:"college"       validate:"required"`
	Github      string `json:"github"        validate:"required"`
	Country     string `json:"country"       validate:"required"`
}
