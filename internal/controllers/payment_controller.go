package controllers

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/parnurzeal/gorequest"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

func generateUniqueRegistrationNo() string {
	timestamp := time.Now().UnixNano()
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(10000000000)
	registrationNo := timestamp*10000000000 + int64(randomNum)
	registrationNo = abs(registrationNo)
	return strconv.FormatInt(registrationNo, 10)
}

func CallBackURL(c *fiber.Ctx) error {
	var request struct {
		RegistrationNo  string  `form:"referenceNo" validate:"required"`
		Token           string  `form:"token" validate:"required"`
		IToken          string  `form:"itoken" validate:"required"`
		TransactionID   string  `form:"transactionId" validate:"required"`
		PaymentStatus   int     `form:"status" validate:"required"`
		Amount          float32 `form:"amount" validate:"required"`
		InvoiceNumber   string  `form:"invoiceNo" validate:"required"`
		TransactionDate string  `form:"transactionDate" validate:"required"`
		// CurrencyCode    string  `form:"currency_code" validate:"required"`
		// UserID          uint    `form:"user_id" validate:"required"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false, "message": "Could not process form data",
		})
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false, "message": "Please pass in all the fields",
		})
	}

	var invoice models.Invoice
	result := database.DB.Where("registration_no = ? ", request.RegistrationNo).First(&invoice)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": false, "message": "Invoice not found"})
	}

	if invoice.PaymentStatus == 2 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": false, "message": "Invoice already paid"})
	}

	if invoice.Amount != request.Amount {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "Invalid details"})
	}

	updates := models.Invoice{
		IToken:          request.IToken,
		Token:           request.Token,
		TransactionId:   request.TransactionID,
		PaymentStatus:   request.PaymentStatus,
		InvoiceNumber:   request.InvoiceNumber,
		TransactionDate: request.TransactionDate,
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
	RegistrationNo := generateUniqueRegistrationNo()

	if country == "IN" {
		amount = float32(1)
		currency_code = "INR"
	}

	invoice := models.Invoice{
		UserID:          user.ID,
		IToken:          "null",
		Token:           "null",
		TransactionId:   "null",
		RegistrationNo:  RegistrationNo,
		PaymentStatus:   0,
		Amount:          amount,
		InvoiceNumber:   "null",
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
	values.Set("registrationNo", RegistrationNo)
	values.Set("amount", fmt.Sprintf("%.2f", amount))
	values.Set("currency_code", currency_code)

	request := gorequest.New().Post("https://events.vit.ac.in/events/bolt/cnfpay")
	request.Type("form")
	request.Send(values.Encode())
	resp, _, errs := request.End()

	if len(errs) > 0 || resp.StatusCode != fiber.StatusOK {
		fmt.Println(errs)
		if err := database.DB.Where("registration_no = ?", RegistrationNo).Unscoped().Delete(&models.Invoice{}).Error; err != nil {
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
		"PaymentLink": "https://events.vit.ac.in/events/bolt/startPay/" + RegistrationNo,
	})
}
