package controllers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"

	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/utils"
)

func Login(c *fiber.Ctx) error {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var user models.User

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{"message": "email or password is missing"})
	}

	database.DB.Find(&user, "email = ?", loginRequest.Email)

	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User does not exist"})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Request",
		})
	}
	user.TokenVersion += 1

	payload := utils.TokenPayload{
		Email:   user.Email,
		Role:    user.Role,
		Version: user.TokenVersion,
	}

	accessToken, err := utils.CreateToken(time.Minute*15, payload, utils.ACCESS_TOKEN, viper.GetString("ACCESS_SECRET_KEY"))
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "Could not sign access token"})
	}
	database.DB.Save(&user)

	refreshToken, err := utils.CreateToken(time.Hour*24, payload, utils.REFRESH_TOKEN, viper.GetString("REFRESH_SECRET_KEY"))
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not sign refresh token",
		})
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "Login Successful", "Access Token": accessToken, "Refresh Token": refreshToken})
}

func Refresh(c *fiber.Ctx) error {
	var tokenReq struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&tokenReq); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "Please pass in a refresh_token"})
	}

	local := c.Locals("user")
	if local == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Please Log In"})
	}

	user := local.(models.User)
	user.TokenVersion += 1

	payload := utils.TokenPayload{
		Email:   user.Email,
		Role:    user.Role,
		Version: user.TokenVersion,
	}

	accessToken, err := utils.CreateToken(time.Minute*15, payload, utils.ACCESS_TOKEN, viper.GetString("ACCESS_SECRET_KEY"))
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "Could not sign access token"})
	}
	database.DB.Save(&user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Access Token": accessToken})
}

func Logout(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	user.TokenVersion = 0

	if user.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User not logged in",
		})
	}

	database.DB.Save(&user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logout Successful"})
}
