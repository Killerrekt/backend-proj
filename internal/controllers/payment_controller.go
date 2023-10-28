package controllers

import (
	"fmt"
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/parnurzeal/gorequest"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

func CallBackURL(c *fiber.Ctx) error {
	var request struct {
		RegistrationNo  int64   `json:"registration_no" validate:"required"`
		Token           string  `json:"token" validate:"required"`
		IToken          string  `json:"itoken" validate:"required"`
		TransactionID   int     `json:"transaction_id" validate:"required"`
		PaymentStatus   bool    `json:"status" validate:"required"`
		Amount          float32 `json:"amount" validate:"required"`
		InvoiceNumber   int64   `json:"invoice_no" validate:"required"`
		TransactionDate string  `json:"transaction_date" validate:"required"`
		CurrencyCode    string  `json:"currency_code" validate:"required"`
		UserID          uint    `json:"user_id" validate:"required"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false, "message": "Could not process JSON",
		})
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false, "message": "Please pass in all the fields",
		})
	}

	var invoice models.Invoice
	result := database.DB.Where("registration_no = ? AND user_id = ?", request.RegistrationNo, request.UserID).First(&invoice)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": false, "message": "Invoice not found"})
	}

	updates := models.Invoice{
		UserID:          request.UserID,
		IToken:          request.IToken,
		Token:           request.Token,
		TransactionId:   request.TransactionID,
		PaymentStatus:   request.PaymentStatus,
		Amount:          request.Amount,
		InvoiceNumber:   request.InvoiceNumber,
		TransactionDate: request.TransactionDate,
		CurrencyCode:    request.CurrencyCode,
	}

	if err := database.DB.Model(&invoice).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false, "message": "Failed to update invoice",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Invoice updated successfully",
	})
}

func InitiatePayment(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var existingInvoice models.Invoice
	if err := database.DB.Where("user_id = ?", user.ID).First(&existingInvoice).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":     "Invoice for this user already exists.",
			"status":      false,
			"PaymentLink": "https://events.vit.ac.in/events/bolt/startPay/" + existingInvoice.RegistrationNo,
		})
	}

	name := user.FirstName + " " + user.LastName
	phoneNumber := user.PhoneNumber
	email := user.Email
	country := user.Country
	currency_code := "USD"
	amount := float32(1)
	RegistrationNo, err := uuid.NewUUID()

	if err != nil {
		fmt.Println("Error generating UUID:", err)
	}

	if country == "IN" {
		amount = float32(1)
		currency_code = "INR"
	}

	invoice := models.Invoice{
		UserID:          user.ID,
		IToken:          "null",
		Token:           "null",
		TransactionId:   0,
		RegistrationNo:  RegistrationNo.String(),
		PaymentStatus:   false,
		Amount:          amount,
		InvoiceNumber:   0,
		TransactionDate: "0",
		CurrencyCode:    currency_code,
	}

	if err := database.DB.Create(&invoice).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false, "message": "Failed to save",
		})
	}

	values := url.Values{}
	values.Set("name", name)
	values.Set("mobileNo", phoneNumber)
	values.Set("email", email)
	values.Set("registrationNo", RegistrationNo.String())
	values.Set("amount", fmt.Sprintf("%.2f", amount))
	values.Set("currency_code", currency_code)

	request := gorequest.New().Post("https://events.vit.ac.in/events/bolt/cnfpay")
	request.Type("form")
	request.Send(values.Encode())
	resp, _, errs := request.End()

	if len(errs) > 0 || resp.StatusCode != fiber.StatusOK {
		fmt.Println(errs)
		if err := database.DB.Where("registration_no = ?", RegistrationNo.String()).Unscoped().Delete(&models.Invoice{}).Error; err != nil {
			fmt.Println("Error deleting entry:", err)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to make the request",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      true,
		"message":     "Invoice generated successfully.",
		"PaymentLink": "https://events.vit.ac.in/events/bolt/startPay/" + RegistrationNo.String(),
	})
}
