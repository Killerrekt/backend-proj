package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/parnurzeal/gorequest"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

func generateUniqueRegistrationNo() int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Intn(9000000) + 1000000)
}

func CallBackURL(c *fiber.Ctx) error {
	var request struct {
		RegistrationNo  int64   `json:"registration_no"`
		Token           string  `json:"token"`
		IToken          string  `json:"itoken"`
		TransactionId   int     `json:"transaction_id"`
		PaymentStatus   bool    `json:"status"`
		Amount          float32 `json:"amount"`
		InvoiceNumber   int64   `json:"invoice_no"`
		TransactionDate string  `json:"transaction_date"`
		CurrencyCode    string  `json:"currency_code"`
		UserId          int     `json:"user_id"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	var invoice models.Invoice
	result := database.DB.Where("registration_no = ? AND user_id = ?", request.RegistrationNo, int(request.UserId)).First(&invoice)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Invoice not found"})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update invoice"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Invoice updated successfully"})
}

func InitiatePayment(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	name := user.FirstName + " " + user.LastName
	phoneNumber := user.PhoneNumber
	email := user.Email
	country := user.Country
	currency_code := "USD"
	amount := float32(20)
	RegistrationNo := generateUniqueRegistrationNo()

	if country == "IN" {
		amount = float32(599)
		currency_code = "INR"
	}

	invoice := models.Invoice{
		UserId:          int(user.ID),
		IToken:          "null",
		TransactionId:   0,
		RegistrationNo:  RegistrationNo,
		PaymentStatus:   false,
		Amount:          amount,
		InvoiceNumber:   0,
		TransactionDate: "0",
		CurrencyCode:    currency_code,
	}

	if err := database.DB.Create(&invoice).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to save"})
	}

	payload := fiber.Map{
		"name":           string(name),
		"mobileNo":       phoneNumber,
		"email":          email,
		"currency_code":  currency_code,
		"amount":         amount,
		"registrationNo": RegistrationNo,
		"formSubmit":     "Submit",
		"mobile_no":      phoneNumber,
	}

	request := gorequest.New().Post("https://events.vit.ac.in/events/GRV23/cnfpay").Send(payload)
	resp, body, errs := request.End()

	if len(errs) > 0 {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to make the request"})
	}

	if resp.StatusCode == http.StatusOK {
		var responseMap map[string]interface{}
		if err := json.Unmarshal([]byte(body), &responseMap); err != nil {
			fmt.Println("Error parsing response:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to parse response"})
		}

		getURL := fmt.Sprintf("https://events.vit.ac.in/events/GRV23/startPay/%d", RegistrationNo)

		getResponse, _, getErrs := gorequest.New().Get(getURL).End()

		if len(getErrs) > 0 {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to make GET request"})
		}

		return c.SendString(fmt.Sprintf("Response from GET request: %s", getResponse))
	}

	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "POST request failed"})
}
