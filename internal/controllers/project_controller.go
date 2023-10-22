package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

func CreateProject(c *fiber.Ctx) error { //this will both create and update the project

	user := c.Locals("user").(models.User)

	var createproject models.Project
	validate := validator.New()

	if err := c.BodyParser(&createproject); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"Error":  "Failed to parse the body",
			"Status": false,
		})
	}

	err := validate.Struct(createproject)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"Error":  "The resquest didn't provide sufficient data",
			"Status": false,
		})
	}

	var team models.Team
	database.DB.Find(&team, "team_id = ?", createproject.TeamID) //maybe changed in future
	if team.TeamID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"Error":  "The team ID provided doesn't exists",
			"Status": false,
		})
	}

	var project models.Project
	database.DB.Find(&project, "team_id = ?", user.TeamId) //maybe changed in future to ID instead of TeamID
	if project.ID != 0 && project.IsFinal {
		return c.Status(fiber.StatusForbidden).JSON(&fiber.Map{
			"Error":  "The project submission have been finalized",
			"Status": false,
		})
	}
	createproject.IsFinal = false

	if project.ID == 0 {
		database.DB.Create(&createproject)
	} else {
		database.DB.Model(&project).Updates(&createproject)
	}
	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"Message": "Route works",
		"Status":  true,
		"Data":    project,
	})
}

func GetProject(c *fiber.Ctx) error {

	var getproject models.Project

	user := c.Locals("user").(models.User)
	database.DB.Find(&getproject, "team_id = ?", user.TeamId)

	if getproject.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"Error":  "The user hasn't created any submission yet that can be viewed",
			"Status": false,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"Data":   getproject,
		"Status": true,
	})
}

func DeleteProject(c *fiber.Ctx) error {
	var deleteproject struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&deleteproject); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"Error":  "Unable to parse the data",
			"Status": false,
		})
	}

	var project models.Project
	database.DB.Find(&project, "name = ?", deleteproject.Name)
	if project.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"Error":  "Project by that name doesn't exists",
			"Status": false,
		})
	}
	database.DB.Delete(&project)
	return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
		"Message": "Successfully deleted the project",
		"Data":    project,
		"Status":  true,
	})
}

func FinaliseProject(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	if !user.IsLeader {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"Error":  "The user is not the leader",
			"Status": false,
		})
	}
	var project models.Project
	database.DB.Find(&project, "team_id = ?", user.TeamId)
	if project.IsFinal {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"Error":  "The project is already finalised",
			"Status": false,
		})
	}
	project.IsFinal = true
	database.DB.Save(&project)
	return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
		"Message": "Project has been finalised",
		"Status":  true,
	})
}

func CreateTeam(c *fiber.Ctx) error { //dummy function just to check functionality
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
			"Error": "Unable to parse the body",
		})
	}
	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"Message": "Team field shld be created",
		"User":    user,
		"Data":    entry,
	})
}
