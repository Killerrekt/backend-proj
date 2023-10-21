package controllers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

func CreateProject(c *fiber.Ctx) error {

	var createproject models.CreateProject
	var user models.User

	validate := validator.New()

	if err := c.BodyParser(&createproject); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"Error": "Unable to parse the req body",
		})
	}

	err := validate.Struct(createproject)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"Error": "The resquest didn't provide sufficient data",
		})
	}

	database.DB.Find(&user, "ID = ?", createproject.UserID)
	if user.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"Error": "User doesn't exist",
		})
	}
	//database.DB.Find() to get the project details of that particular user through team table

	var project models.Project
	database.DB.Find(&project, "ID = ?", 6) //here the ID will the one from the team table
	fmt.Printf("project is %d", project.ID)
	if project.ID != 0 && project.IsFinal {
		return c.Status(fiber.StatusForbidden).JSON(&fiber.Map{
			"Error": "The project submission have been finalized",
		})
	}

	proj := models.Project{
		Name:         createproject.Name,
		Desc:         createproject.Desc,
		Githublink:   createproject.Githublink,
		FigmaLink:    createproject.FigmaLink,
		VideoLink:    createproject.VideoLink,
		DriveLink:    createproject.DriveLink,
		ProjectTrack: createproject.ProjectTrack,
		IsFinal:      false,
	}

	if project.ID == 0 {
		database.DB.Create(&proj)
	} else {
		database.DB.Model(&project).Where("ID = ?", project.ID).Updates(&proj)
	}
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

	return nil
}

func CreateTeam(c *fiber.Ctx) error {
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
	})
}
