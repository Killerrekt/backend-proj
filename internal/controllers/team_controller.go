package controllers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/services"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/utils"
)

func CreateTeam(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var data struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "error parsing JSON",
		})
	}

	if data.Name == "" {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"status":  false,
			"message": "name is required",
		})
	}

	if user.TeamID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "already in team",
		})
	}

	_, err := services.FindTeamByName(data.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": false, "message": "Some error occurred", "error": err.Error()})
	}

	code, err := utils.GenerateUniqueTeamCode()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": false, "message": "Some error occurred", "error": err.Error()})
	}

	team := models.Team{
		Name:     data.Name,
		TeamID:   uint(uuid.New().ID()),
		LeaderID: user.ID,
		Code:     code,
	}

	if err := database.DB.Create(&team).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Could not create team",
		})
	}

	user.IsLeader = true
	user.TeamID = team.TeamID

	if err := database.DB.Save(user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Could not set user as leader",
		})
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"status": true, "code": code, "message": "Team created successfully"})
}

func JoinTeam(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var data struct {
		Code string `json:"code"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false, "message": "error parsing JSON",
		})
	}

	if user.TeamID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "already in team",
		})
	}

	team, err := services.FindTeamByCode(data.Code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).
				JSON(fiber.Map{"status": false, "message": "Team does not exist"})
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": false, "message": "Some error occurred", "error": err.Error()})
	}

	if len(team.Users) >= 4 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "team is full",
		})
	}
	user.TeamID = team.TeamID

	if err := database.DB.Save(&team).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Could not create team",
		})
	}

	if err := database.DB.Save(&user); err.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  false,
			"message": "could not save user team id",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "message": "Joined team"})
}

func GetTeam(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	id := user.TeamID

	var team models.Team
	if err := database.DB.Preload(clause.Associations).First(&team, "team_id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "team not found",
			"data":    team,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "team": team})
}

func UpdateTeam(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	id := user.TeamID

	var data struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "error parsing JSON",
		})
	}

	if data.Name == "" {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"status":  false,
			"message": "name is required",
		})
	}

	var team models.Team
	if err := database.DB.Find(&team, "team_id = ?", id); err.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "team not found",
		})
	}

	team.Name = data.Name

	if err := database.DB.Save(&team).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to save team",
		})
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"status": true, "message": "Updated successfully"})
}

func DeleteTeam(c *fiber.Ctx) error {

	user := c.Locals("user").(models.User)
	id := user.TeamID

	/*id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "please give a valid ID"})
	}*/

	/*err := */
	if !user.IsLeader {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "user is not a leader",
			"status": false,
		})
	}
	services.DeleteTeamByID(uint(id))

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "team deleted",
	})
}

func GetAllTeams(c *fiber.Ctx) error {
	var teams []models.Team

	if err := database.DB.Preload(clause.Associations).Find(&teams).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to get teams",
		})
	}

	return c.JSON(teams)
}

func GetLeaderInfo(c *fiber.Ctx) error {
	id := c.Params("id")

	var user models.User
	if err := database.DB.Find(&user, "team_id = ?", id).Where("is_leader = true").Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Leader not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "user": user})
}

func GetProjectFromTeamID(c *fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "Please pass in a valid id"})
	}

	team, err := services.FindTeamByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": false, "message": "Team does not exist"})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "project": team.Project})
}

func GetIdeaFromTeamID(c *fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "Please pass in a valid id"})
	}

	team, err := services.FindTeamByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": false, "message": "Team does not exist"})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "project": team.Idea})
}

func LeaveTeam(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	if user.TeamID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "User not part of any team"})
	}

	user.TeamID = 0
	database.DB.Save(&user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "message": "Left team"})
}
