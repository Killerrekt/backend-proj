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

func UpdateIdea(c *fiber.Ctx) error {
	var req struct {
		FigmaLink string `json:"figma_link"`
		DriveLink string `json:"drive_link"`
		VideoLink string `json:"video_link"`
		Desc      string `json:"desc"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "unable to parse the body",
			"status": false,
		})
	}

	user := c.Locals("user").(models.User)
	var exists models.Idea
	database.DB.Find(&exists, "team_id = ?", user.TeamID)
	if exists.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "No idea exist under this user(try create route)",
			"status": false,
		})
	}

	if req.FigmaLink != "" {
		exists.FigmaLink = req.FigmaLink
	}
	if req.DriveLink != "" {
		exists.DriveLink = req.DriveLink
	}
	if req.VideoLink != "" {
		exists.VideoLink = req.VideoLink
	}
	if req.Desc != "" {
		exists.Desc = req.Desc
	}

	if err := database.DB.Save(&exists); err.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "something went wrong in updating",
			"status": false,
		})
	}
	return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
		"message": "idea updated",
		"data":    exists,
		"status":  true,
	})
}

func DeleteIdea(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	var exist models.Idea
	database.DB.Find(&exist, "team_id = ?", user.TeamID)
	if exist.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "No idea exists that can be deleted",
			"status": false,
		})
	}
	if err := database.DB.Unscoped().Delete(&exist); err.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "something went wrong while deleting the idea",
			"status": false,
		})
	}
	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"message": "Idea is deleted",
		"status":  true,
	})
}

func GetIdea(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	var exist models.Idea

	database.DB.Find(&exist, "team_id = ?", user.TeamID)
	if exist.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "No idea exists",
			"status": false,
		})
	}
	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"message": "successfully fetched the data",
		"data":    exist,
		"status":  true,
	})
}

func GetAllIdea(c *fiber.Ctx) error {
	var allideas []models.Idea
	database.DB.Find(&allideas)
	if len(allideas) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "No idea exists",
			"status":  true,
		})
	}
	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"data":   allideas,
		"status": true,
	})
}
