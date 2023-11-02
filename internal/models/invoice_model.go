package models

import "gorm.io/gorm"

type Invoice struct {
	gorm.Model
	UserID          uint    `json:"user_id"`
	IToken          string  `json:"itoken"`
	Token           string  `json:"token"`
	TransactionId   string  `json:"transaction_id"`
	RegistrationNo  string  `json:"registration_no" gorm:"unique"`
	PaymentStatus   int     `json:"status"`
	Amount          float32 `json:"amount"`
	InvoiceNumber   string  `json:"invoice_no"`
	TransactionDate string  `json:"transaction_date"`
	CurrencyCode    string  `json:"currency_code"`
}
