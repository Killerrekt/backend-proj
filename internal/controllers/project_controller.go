package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

func CreateProject(c *fiber.Ctx) error { // this will both create and update the project

	user := c.Locals("user").(models.User)

	var createproject models.Project
	validate := validator.New()

	if err := c.BodyParser(&createproject); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": false,
			"error":  "Failed to parse the body",
		})
	}

	err := validate.Struct(createproject)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": false,
			"error":  "The resquest didn't provide sufficient data",
		})
	}

	var team models.Team
	database.DB.Find(&team, "team_id = ?", createproject.TeamID) // maybe changed in future
	if team.TeamID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": false,
			"error":  "The team ID provided doesn't exists",
		})
	}

	var project models.Project
	database.DB.Find(&project, "team_id = ?", user.TeamId) // maybe changed in future to ID instead of TeamID
	if project.ID != 0 && project.IsFinal {
		return c.Status(fiber.StatusForbidden).JSON(&fiber.Map{
			"status": false,
			"error":  "The project submission have been finalized",
		})
	}
	createproject.IsFinal = false

	if project.ID == 0 {
		database.DB.Create(&createproject)
	} else {
		database.DB.Model(&project).Updates(&createproject)
	}
	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"status":  true,
		"message": "Route works",
		"data":    project,
	})
}

func GetProject(c *fiber.Ctx) error {
	var getproject models.Project

	user := c.Locals("user").(models.User)
	database.DB.Find(&getproject, "team_id = ?", user.TeamId)

	if getproject.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": false,
			"error":  "The user hasn't created any submission yet that can be viewed",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"status": true,
		"data":   getproject,
	})
}

func DeleteProject(c *fiber.Ctx) error {
	var deleteproject struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&deleteproject); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": false,
			"error":  "Unable to parse the data",
		})
	}

	var project models.Project
	database.DB.Find(&project, "name = ?", deleteproject.Name)
	if project.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": false,
			"error":  "Project by that name doesn't exists",
		})
	}
	database.DB.Delete(&project)
	return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
		"status":  true,
		"message": "Successfully deleted the project",
		"data":    project,
	})
}

func FinaliseProject(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	if !user.IsLeader {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": false,
			"error":  "The user is not the leader",
		})
	}
	var project models.Project
	database.DB.Find(&project, "team_id = ?", user.TeamId)
	if project.IsFinal {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": false,
			"error":  "The project is already finalised",
		})
	}
	project.IsFinal = true
	database.DB.Save(&project)
	return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
		"status":  true,
		"message": "Project has been finalised",
	})
}

func CreateTeam(c *fiber.Ctx) error { // dummy function just to check functionality
	user := c.Locals("user").(models.User)
	var Req struct {
		TeamID int `json:"team_id"`
	}

	err := c.BodyParser(&Req)

	user.TeamId = Req.TeamID
	database.DB.Save(&user)

	entry := models.Team{
		TeamID: uint(Req.TeamID),
	}

	database.DB.Create(&entry)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": false,
			"error":  "Unable to parse the body",
		})
	}
	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"status":  true,
		"message": "Team field shld be created",
		"user":    user,
		"data":    entry,
	})
}