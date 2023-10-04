package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
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

func ForgotEmail(c *fiber.Ctx) error {

	type EmailData struct {
		Email string `json:"email"`
	}

	var email EmailData

	if err := c.BodyParser(&email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"Status": false,
			"Error":  err,
		})
	}

	var check models.User

	database.DB.Find(&check, "email = ?", email.Email)

	if (check == models.User{}) {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"Status": false,
			"Error":  "The email address given doesn't exist",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":   time.Now().Add(time.Minute * 2).Unix(),
		"email": email.Email,
	})

	Token, err := token.SignedString([]byte(viper.GetString("ACCESS_SECRET_KEY")))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"Status": false,
			"Error":  "Failed to create an JWT token",
		})
	}

	//msg := gomail.NewMessage()
	//msg.SetHeader("FROM",)

	//token has been created and stored in Token

	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"data": Token,
	})
}
