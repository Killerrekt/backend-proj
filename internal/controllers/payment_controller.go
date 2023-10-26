package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/parnurzeal/gorequest"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

func generateUniqueRegistrationNo() int64 {
	return int64(rand.Intn(9000000) + 1000000)
}

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
		TransactionID:   request.TransactionID,
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
	currencyCode := "USD"
	amount := float32(20)
	RegistrationNo := generateUniqueRegistrationNo()

	if country == "IN" {
		amount = float32(599)
		currencyCode = "INR"
	}

	invoice := models.Invoice{
		UserID:          user.ID,
		IToken:          "null",
		TransactionID:   0,
		RegistrationNo:  RegistrationNo,
		PaymentStatus:   false,
		Amount:          amount,
		InvoiceNumber:   0,
		TransactionDate: "0",
		CurrencyCode:    currencyCode,
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
		"currency_code":  currencyCode,
		"amount":         amount,
		"registrationNo": RegistrationNo,
		"formSubmit":     "Submit",
		"mobile_no":      phoneNumber,
	}

	request := gorequest.New().Post("https://events.vit.ac.in/events/GRV23/cnfpay").Send(payload)
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

		getURL := fmt.Sprintf("https://events.vit.ac.in/events/GRV23/startPay/%d", RegistrationNo)

		getResponse, _, getErrs := gorequest.New().Get(getURL).End()

		if len(getErrs) > 0 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": false, "message": "Failed to make GET request",
			})
		}

		return c.SendString(fmt.Sprintf("Response from GET request: %v", getResponse))
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status":  true,
		"message": "POST request failed",
	})
}
