package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

func CreateIdea(c *fiber.Ctx) error {
	var req models.Idea

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "unable to parse the body",
			"status": false,
		})
	}

	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "Missing required fields(Title and Description)",
			"status": false,
		})
	}

	user := c.Locals("user").(models.User)
	var check models.Team
	database.DB.Find(&check, "team_id = ?", user.TeamID)
	if check.TeamID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "Team ID of the user doesn't exists",
			"status": false,
		})
	}

	var exists models.Idea
	database.DB.Find(&exists, "team_id = ?", user.TeamID)
	if exists.ID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "Idea already exists(use update route)",
			"status": false,
		})
	}

	req.TeamID = user.TeamID
	if err := database.DB.Create(&req); err.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":    "Something went wrong while saving the idea",
			"db_error": err,
			"status":   false,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"message": "idea successfully created",
		"data":    req,
		"status":  true,
	})
}
