package models

import "gorm.io/gorm"

type Invoice struct {
	gorm.Model
	UserId          int     `json:"user_id" gorm:"unique"`
	IToken          string  `json:"itoken"`
	TransactionId   int     `json:"transaction_id"`
	ReferenceNo     int64   `json:"reference_no" gorm:"unique"`
	RegistrationNo  int     `json:"registration_no" gorm:"unique"`
	PaymentStatus   bool    `json:"status"`
	Amount          float32 `json:"amount"`
	InvoiceNumber   int64   `json:"invoice_no"`
	TransactionDate string  `json:"transaction_date"`
	CurrencyCode    string  `json:"currency_code"`
}
