package controllers

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

func CreateProject(c *fiber.Ctx) error { // this will both create

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
			"error":  "The request didn't provide sufficient data",
		})
	}

	var team models.Team
	database.DB.Find(&team, "team_id = ?", user.TeamID) // maybe changed in future
	if team.TeamID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": false,
			"error":  "The team ID provided doesn't exists",
		})
	}

	var project models.Project
	database.DB.Find(&project, "team_id = ?", user.TeamID) // maybe changed in future to ID instead of TeamID
	if project.ID != 0 && project.IsFinal {
		return c.Status(fiber.StatusForbidden).JSON(&fiber.Map{
			"status": false,
			"error":  "The project submission have been finalized",
		})
	}

	fmt.Println(project)

	if project.ID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "the project is already created",
			"status": false,
		})
	}
	createproject.IsFinal = false

	var errstring string

	if project.ID == 0 {
		err := database.DB.Create(&createproject)
		errstring = DBerrorHandling(err)
	} else {
		err := database.DB.Model(&project).Updates(&createproject)
		errstring = DBerrorHandling(err)
	}
	if errstring != "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  errstring,
			"status": false,
		})
	}
	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"status":  true,
		"message": "project have been created",
		"data":    project,
	})
}

func GetProject(c *fiber.Ctx) error {
	var getproject models.Project

	user := c.Locals("user").(models.User)
	err := database.DB.Find(&getproject, "team_id = ?", user.TeamID)
	errstring := DBerrorHandling(err)
	if errstring != "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  errstring,
			"status": false,
		})
	}

	if getproject.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": false,
			"error":  "The user hasn't created any project yet that can be viewed",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"status": true,
		"data":   getproject,
	})
}

func GetAllProject(c *fiber.Ctx) error {
	var data []models.Project
	err := database.DB.Find(&data)
	errstring := DBerrorHandling(err)
	if errstring != "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  errstring,
			"status": false,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"data":   data,
		"status": true,
	})
}

func DeleteProject(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	var project models.Project
	database.DB.Find(&project, "team_id = ?", user.TeamID)
	if project.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": false,
			"error":  "Project by the user doesn't exists",
		})
	}
	err := database.DB.Unscoped().Delete(&project)
	if check := DBerrorHandling(err); check != "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  check,
			"status": false,
		})
	}
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
	database.DB.Find(&project, "team_id = ?", user.TeamID)
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

func UpdateProject(c *fiber.Ctx) error {
	var updateproject models.CreateProject

	user := c.Locals("user").(models.User)

	if err := c.BodyParser(&updateproject); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "failed to parse the body",
			"status": false,
		})
	}

	var currproject models.Project
	err := database.DB.Find(&currproject, "team_id = ?", user.TeamID)
	if check := DBerrorHandling(err); check != "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  check,
			"status": false,
		})
	}

	if currproject.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "no project exist for this team",
			"status": false,
		})
	} else if currproject.IsFinal {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "the project has been finalized and can't be changed",
			"status": true,
		})
	}

	if updateproject.Desc != "" {
		currproject.Desc = updateproject.Desc
	}
	if updateproject.DriveLink != "" {
		currproject.DriveLink = updateproject.DriveLink
	}
	if updateproject.FigmaLink != "" {
		currproject.FigmaLink = updateproject.FigmaLink
	}
	if updateproject.Githublink != "" {
		currproject.Githublink = updateproject.Githublink
	}
	if updateproject.VideoLink != "" {
		currproject.VideoLink = updateproject.VideoLink
	}

	err = database.DB.Save(&currproject)
	if check := DBerrorHandling(err); check != "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  check,
			"status": false,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"data":   currproject,
		"status": true,
	})
}

/*func CreateTeam(c *fiber.Ctx) error { // dummy function just to check functionality
	user := c.Locals("user").(models.User)
	var Req struct {
		TeamID uint `json:"team_id"`
	}

	err := c.BodyParser(&Req)

	user.TeamID = Req.TeamID
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
		"message": "Team field should be created",
		"user":    user,
		"data":    entry,
	})
}*/

func DBerrorHandling(err *gorm.DB) string {
	if err.Error != nil {
		if errors.Is(err.Error, gorm.ErrDuplicatedKey) {
			return "Duplicate field was tried to be entered"
		} else {
			return err.Error.Error()
		}
	}
	return ""
}
