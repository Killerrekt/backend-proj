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
	// Get user ID and team name from request parameters
	userID := c.Params("user_id")
	teamName := c.Params("team_name")

	// Fetch user using user_id
	user, err := findUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	// Check if user is already part of a team
	if user.TeamID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user already in a team"})
	}

	// Check if team name already exists
	if _, err := findTeamByName(teamName); err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "team name already taken"})
	}

	// Create team
	team := models.Team{
		Name:     teamName,
		TeamCode: uuid.New().String(),
	}

	// Save team to get ID
	if err := database.DB.Create(&team).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create team"})
	}

	// Get team ID
	teamID := team.ID

	// Assign team to user
	user.TeamID = teamID
	user.IsLeader = true

	// Save user
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save user"})
	}

	// Reload team to get associations
	if err := database.DB.Preload("Members").First(&team, teamID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to load team"})
	}

	return c.JSON(team)
}

// JOIN TEAM
func JoinTeam(c *fiber.Ctx) error {
	// Parse request parameters
	var data struct {
		UserID string `json:"user_id"`
		Code   string `json:"code"`
	}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// Find user
	user, err := findUserByID(data.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	// Check if user is already part of a team
	if user.TeamID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user already in a team"})
	}

	// Find team by code or ID
	team, err := findTeamByCodeOrID(data.Code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	// Check if team is full
	if len(team.Members) >= MaxTeamMembers {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "team full"})
	}

	// Add user to team
	user.TeamID = team.ID
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save user"})
	}

	return c.JSON(team)
}

// GET TEAM
func GetTeam(c *fiber.Ctx) error {
	// Get team ID
	id := c.Params("id")

	// Find team
	team, err := findTeamByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	// Set members count
	team.MembersCount = len(team.Members)

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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// Find team
	team, err := findTeamByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	// Update name
	team.Name = data.Name

	// Save team
	if err := database.DB.Save(&team).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update team"})
	}

	return c.JSON(team)
}

// DELETE TEAM
func DeleteTeam(c *fiber.Ctx) error {
	// Get team ID
	id := c.Params("id")

	// Find team
	team, err := findTeamByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	// Delete team
	if err := database.DB.Delete(&team).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete team"})
	}

	// Remove team from members
	if err := removeTeamFromMembers(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to remove team from members"})
	}

	return c.JSON(fiber.Map{"message": "team deleted"})
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
