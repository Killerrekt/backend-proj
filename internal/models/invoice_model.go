package models

import "gorm.io/gorm"

type Invoice struct {
	gorm.Model
	UserId          int     `json:"user_id"`
	IToken          string  `json:"itoken"`
	Token           string  `json:"token"`
	TransactionId   int     `json:"transaction_id"`
	RegistrationNo  string  `json:"registration_no" gorm:"unique"`
	PaymentStatus   bool    `json:"status"`
	Amount          float32 `json:"amount"`
	InvoiceNumber   int64   `json:"invoice_no"`
	TransactionDate string  `json:"transaction_date"`
	CurrencyCode    string  `json:"currency_code"`
}
