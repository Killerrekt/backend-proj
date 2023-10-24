package controllers

import (
	"encoding/json"
	"fmt"

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
		TransactionId   int     `json:"transaction_id" validate:"required"`
		PaymentStatus   bool    `json:"status" validate:"required"`
		Amount          float32 `json:"amount" validate:"required"`
		InvoiceNumber   int64   `json:"invoice_no" validate:"required"`
		TransactionDate string  `json:"transaction_date" validate:"required"`
		CurrencyCode    string  `json:"currency_code" validate:"required"`
		UserId          int     `json:"user_id" validate:"required"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false, "message": "Could not process JSON",
		})
	}

	validator := validator.New()

	if err := validator.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false, "message": "Please pass in all the fields",
		})
	}

	var invoice models.Invoice
	result := database.DB.Where("registration_no = ? AND user_id = ?", request.RegistrationNo, int(request.UserId)).First(&invoice)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": false, "message": "Invoice not found"})
	}

	updates := models.Invoice{
		UserId:          request.UserId,
		IToken:          request.IToken,
		Token:           request.Token,
		TransactionId:   request.TransactionId,
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
		UserId:          int(user.ID),
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

	payload := fiber.Map{
		"name":           string(name),
		"mobileNo":       phoneNumber,
		"email":          email,
		"registrationNo": RegistrationNo,
		"amount":         amount,
		"currency_code":  currency_code,
	}

	request := gorequest.New().Post("https://events.vit.ac.in/events/technext/cnfpay").Send(payload)
	resp, body, errs := request.End()

	if len(errs) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false, "message": "Failed to make the request",
		})
	}

	if resp.StatusCode == fiber.StatusOK {
		var responseMap map[string]interface{}
		if err := json.Unmarshal([]byte(body), &responseMap); err != nil {
			fmt.Println("Error parsing response:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": false, "message": "Failed to parse response",
			})
		}
		redirectURL := "https://events.vit.ac.in/events/technext/startPay/" + RegistrationNo.String()
		return c.Redirect(redirectURL, fiber.StatusFound)
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status":  true,
		"message": "POST request failed",
	})
}
