package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

const MaxTeamMembers = 4

// CREATE TEAM
func CreateTeam(c *fiber.Ctx) error {
	// Get logged in user
	user := c.Locals("user").(*models.User)

	// Validate request
	var data struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Check if user already in team
	if user.TeamID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "already in team",
		})
	}

	// Check if name exists
	var existing models.Team
	if err := database.DB.Where("name = ?", data.Name).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name exists",
		})
	}

	// Create team
	team := models.Team{
		Name:   data.Name,
		TeamID: uint(uuid.New().ID()), // Generate a hashed UUID
	}

	// Add logged in user
	team.Users = append(team.Users, *user)

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
	user := c.Locals("user").(*models.User)

	// Validate request
	var data struct {
		Code string `json:"code"`
	}

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Check if user already in team
	if user.TeamID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "already in team",
		})
	}

	// Find team
	var team models.Team
	if err := database.DB.Where("team_code = ?", data.Code).Preload("Users").First(&team).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "invalid code",
		})
	}

	// Check if full
	if len(team.Users) >= 4 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "team full",
		})
	}

	// Add user
	team.Users = append(team.Users, *user)

	// Save
	if err := database.DB.Save(&team).Error; err != nil {
		return err
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

	// Get logged in user
	user := c.Locals("user").(*models.User)

	// Set custom fields
	team.MembersCount = len(team.Users)

	// Check if the user is a leader
	for _, u := range team.Users {
		if u.ID == user.ID && u.IsLeader {
			team.IsUserLeader = true
			break
		}
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
		return err
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
			"error": "team not found",
		})
	}

	// Delete team
	if err := database.DB.Delete(&team).Error; err != nil {
		return err
	}

	// Remove team from members
	if err := database.DB.Model(&models.User{}).Where("team_id = ?", id).Update("team_id", 0).Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
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
