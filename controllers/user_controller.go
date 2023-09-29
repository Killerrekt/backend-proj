package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

func CreateUser(c *fiber.Ctx) error {
	var createUser models.CreateUser

	if err := c.BodyParser(&createUser); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "Please send complete data"})
	}

	dob, _ := time.Parse("2006-01-02", createUser.DateOfBirth)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(createUser.Password), 10)

	user := models.User{
		FirstName:   createUser.FirstName,
		LastName:    createUser.LastName,
		Email:       createUser.Email,
		Password:    string(hashedPassword),
		Gender:      createUser.Gender,
		DateOfBirth: dob,
		Bio:         createUser.Bio,
		TeamId:      0,
		IsLeader:    false,
		IsApproved:  false,
		PhoneNumber: createUser.PhoneNumber,
		College:     createUser.College,
		Github:      createUser.Github,
	}

	if result := database.DB.Create(&user); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": result.Error.Error()})
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "Successfully created user", "user": user})
}
