package models

import "gorm.io/gorm"

type Enquiry struct {
	gorm.Model
	Name        string `json:"name" validate:"required"`
	JobProfile  string `json:"job_profile" validate:"required"`
	CompanyName string `json:"company_name" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	Email       string `json:"email" validate:"required"`
	City        string `json:"city" validate:"required"`
	Message     string `json:"message"`
}
