package controllers

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"

	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/utils"
)

func Login(c *fiber.Ctx) error {
	var loginRequest struct {
		Email    string `json:"email"    validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	var user models.User

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{"message": "Could not parse JSON"})
	}

	validator := validator.New()

	if err := validator.Struct(loginRequest); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{"message": "Please pass in all the fields"})
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

	payload := utils.TokenPayload{
		Email: user.Email,
		Role:  user.Role,
	}

	refreshToken, err := utils.CreateToken(
		time.Hour*24,
		payload,
		utils.REFRESH_TOKEN,
		viper.GetString("REFRESH_SECRET_KEY"),
	)
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not sign refresh token",
		})
	}

	if err := database.RedisClient.Set(refreshToken, user.Email, time.Hour*24); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "Some error occured"})
	}

	log.Println("Set Refresh token successful")

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "Login Successful", "Access Toekn": "To get an access token please send a request to the refresh route", "Refresh Token": refreshToken})
}

func Refresh(c *fiber.Ctx) error {
	var tokenReq struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.BodyParser(&tokenReq); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "Could not process JSON"})
	}

	validator := validator.New()

	if err := validator.Struct(tokenReq); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{"message": "Please pass in the correct data"})
	}

	local := c.Locals("user")

	user := local.(models.User)
	payload := utils.TokenPayload{
		Email: user.Email,
		Role:  user.Role,
	}

	accessToken, err := utils.CreateToken(
		time.Minute*15,
		payload,
		utils.ACCESS_TOKEN,
		viper.GetString("ACCESS_SECRET_KEY"),
	)
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "Could not sign access token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Access Token": accessToken})
}

func Logout(c *fiber.Ctx) error {
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error parsing JSON"})
	}

	if request.RefreshToken == "" {
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{"message": "Please pass in the refresh token"})
	}

  if _, err := database.RedisClient.Get(request.RefreshToken); err != nil {
    if err == redis.Nil {
      return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User not logged in"})
    }
  }

	if err := database.RedisClient.Delete(request.RefreshToken); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "Some error occured"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logout Successful"})
}
