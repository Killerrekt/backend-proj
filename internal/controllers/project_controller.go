package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

func CreateProject(c *fiber.Ctx) error {
	var createproject models.CreateProject
	if err := c.BodyParser(&createproject); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"Error": "Unable to parse the req body",
		})
	}

	database.DB.Save(&createproject)

	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"Message": "Route works",
	})
}

func GetProject(c *fiber.Ctx) error {
	type GetProject struct {
		UserID string `json:"userID"`
		TeamID string `json:"teamID"`
	}

	var req GetProject

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"Error": "unable to parse the data",
		})
	}

	fmt.Println(req)

	return nil
}
