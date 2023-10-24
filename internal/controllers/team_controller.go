package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

const MaxTeamMembers = 4

// CREATE TEAM
func CreateTeam(c *fiber.Ctx) error {
	// Get logged in user
	user := c.Locals("user").(models.User)

	// Validate request
	var data struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false, "message": "error parsing JSON",
		})
	}

	if data.Name == "" {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"status":  false,
			"message": "name is required",
		})
	}

	// Check if user already in team
	if user.TeamID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "already in team",
		})
	}

	// Check if name exists
	var existing models.Team
	if err := database.DB.Where("name = ?", data.Name).First(&existing).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"error":  "A team with the same name already exists",
		})
	}

	// Create team
	team := models.Team{
		Name:     data.Name,
		TeamID:   uint(uuid.New().ID()), // Generate a hashed UUID
		LeaderID: user.ID,
	}

	// Save
	if err := database.DB.Create(&team).Error; err != nil {
		return err
	}

	// Set as leader
	user.IsLeader = true

	if err := database.DB.Save(user).Error; err != nil {
		return err
	}

	return c.JSON(team)
}

// JOIN TEAM
func JoinTeam(c *fiber.Ctx) error {
	// Get logged in user
	user := c.Locals("user").(models.User)

	// Validate request
	var data struct {
		Code string `json:"code"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false, "message": "error parsing JSON",
		})
	}

	// Check if user already in team
	if user.TeamID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "already in team",
		})
	}

	// Find team
	var team models.Team
	if err := database.DB.Where("team_code = ?", data.Code).Preload("Users").First(&team).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "invalid code",
		})
	}

	// Check if full
	if team.MembersCount >= 4 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "team is full",
		})
	}

	// Save
	if err := database.DB.Save(&team).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Could not create team",
		})
	}

	return c.JSON(team)
}

// GET TEAM
func GetTeam(c *fiber.Ctx) error {
	// Get team ID
	id := c.Params("id")

	// Find team
	var team models.Team
	if err := database.DB.Preload("Users").First(&team, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "team not found",
		})
	}

	return c.JSON(team)
}

// UPDATE TEAM
func UpdateTeam(c *fiber.Ctx) error {
	// Get team ID
	id := c.Params("id")

	// Validate request
	var data struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false, "message": "error parsing JSON",
		})
	}

	if data.Name == "" {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"status":  false,
			"message": "name is required",
		})
	}

	// Find team
	var team models.Team
	if err := database.DB.First(&team, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "team not found",
		})
	}

	// Update name
	team.Name = data.Name

	// Save
	if err := database.DB.Save(&team).Error; err != nil {
		return err
	}

	return c.JSON(team)
}

// DELETE TEAM
func DeleteTeam(c *fiber.Ctx) error {
	// Get team ID
	id := c.Params("id")

	// Find team
	var team models.Team
	if err := database.DB.First(&team, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"error":  "team not found",
		})
	}

	// Delete team
	if err := database.DB.Delete(&team).Error; err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Some error occured check logs",
		})
	}

	// Remove team from members
	if err := database.DB.Model(&models.User{}).Where("team_id = ?", id).Update("team_id", 0).Error; err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Some error occured check logs",
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "team deleted",
	})
}

// Helper function to check if the user is a leader in the team
func isUserLeader(user *models.User, users []models.User) bool {
	for _, u := range users {
		if u.ID == user.ID && u.IsLeader {
			return true
		}
	}
	return false
}

// Find user by ID
func findUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Find team by name
func findTeamByName(name string) (*models.Team, error) {
	var team models.Team
	if err := database.DB.Where("name = ?", name).First(&team).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

// Find team by code or ID
func findTeamByCodeOrID(codeOrID string) (*models.Team, error) {
	var team models.Team
	if err := database.DB.Where("team_code = ? OR id = ?", codeOrID, codeOrID).First(&team).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

// Find team by ID
func findTeamByID(id string) (*models.Team, error) {
	var team models.Team
	if err := database.DB.Preload("Members").First(&team, id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

// Remove team from members
func removeTeamFromMembers(teamID string) error {
	return database.DB.Model(&models.User{}).Where("team_id = ?", teamID).Update("team_id", 0).Error
}

// GET ALL TEAMS
// GET LEADER INFO
