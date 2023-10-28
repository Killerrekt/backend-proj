package controllers

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

func ExhibitionEnquiry(c *fiber.Ctx) error {
	var request models.Enquiry

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false, "message": "Error parsing JSON",
		})
	}

	validator := validator.New()
	if err := validator.Struct(request); err != nil {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"status": false, "message": "Please pass in all the required fields",
		})
	}

	result := database.DB.Create(&request)
	if result.Error != nil {
		log.Println(result.Error.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false, "message": "Error storing record",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "message": "Response Recorded"})
}
