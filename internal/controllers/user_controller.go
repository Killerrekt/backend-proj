package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"

	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/utils"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "Successfully created user", "user": user})
}

func ForgotPassword(c *fiber.Ctx) error {
	email := c.Params("email")

	var check models.User
	database.DB.Find(&check, "email = ?", email)
	if (check == models.User{}) {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"Status": false,
			"Error":  "The email address given doesn't exist",
		})
	}

	payload := utils.TokenPayload{
		Email:   email,
		Role:    "",
		Version: 0,
	}

	resetToken, err := utils.CreateToken(time.Minute*2, payload, utils.REFRESH_TOKEN, viper.GetString("RESET_SECRET_KEY"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"Status": false,
			"Error":  "Failed to create an JWT token",
		})
	}

	url := fmt.Sprintf("%s%s", viper.GetString("URL"), resetToken)
	message := fmt.Sprintf("%s\n%s %s\n%s",
		"This is an auto generated email.",
		"Click the link below to reset your password",
		url,
		"If this request was not sent by you please report to the concerned authorities")

	err = utils.SendMail("Password Reset", email, message)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"Status": false,
			"Error":  "Something went wrong while sending the email",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"Status": true,
		"data":   resetToken,
	})
}

func ResetPassword(c *fiber.Ctx) error {
	token := c.Params("Token")

	type Password struct {
		Password     string `json:"password"`
		Confirm_pass string `json:"confirm_pass"`
	}

	Token, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(viper.GetString("RESET_SECRET_KEY")), nil
	})

	if decoded, ok := Token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > decoded["exp"].(float64) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Error": "Token Expired",
			})
		}

		email := decoded["email"]
		var user models.User
		database.DB.Find(&user, "email = ?", email)
		if (user == models.User{}) {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"Error": "The email doesn't exists",
			})
		}

		req := new(Password)
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"Error": "Error rose while parsing through the body",
			})
		}

		if req.Password != req.Confirm_pass {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"Error": "Password and confirm password are not the same",
			})
		}

		hashed_password, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		user.Password = string(hashed_password)
		database.DB.Save(user)
		return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
			"Message": "The password has been updated",
		})
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Token"})
}
