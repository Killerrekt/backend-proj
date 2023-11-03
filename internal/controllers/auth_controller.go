package controllers

import (
	"errors"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/services"
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
			JSON(fiber.Map{"status": false, "message": "Could not parse JSON"})
	}

	validator := validator.New()

	if err := validator.Struct(loginRequest); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{"status": false, "message": "Please pass in all the fields"})
	}

	database.DB.Find(&user, "email = ?", loginRequest.Email)

	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"status": false, "message": "User does not exist"})
	}

	if user.IsBanned {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false, "message": "User is banned",
			"verification_status": user.IsVerified, "banned": true,
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status": false, "message": "invalid password",
			"verification_status": user.IsVerified, "banned": false,
		})
	}

	if !user.IsVerified {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{
				"status": false, "message": "User is not verified",
				"verification_status": false, "banned": false,
			})
	}

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
			JSON(fiber.Map{
				"status": false, "message": "Could not sign access token",
			})
	}

	refreshToken, err := utils.CreateToken(
		time.Hour*24,
		payload,
		utils.REFRESH_TOKEN,
		viper.GetString("REFRESH_SECRET_KEY"),
	)
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": false, "message": "Could not sign refresh token"})
	}

	if err := database.RedisClient.Set(refreshToken, user.Email, time.Hour*24); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": false, "message": "Some error occured"})
	}

	inTeam := false
	if user.TeamID != 0 {
		inTeam = true
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{
			"status": true, "message": "Login Successful", "access_token": accessToken,
			"refresh_token": refreshToken, "verification_status": user.IsVerified,
			"payment_status": user.IsPaid, "banned": false, "in_team": inTeam,
		})
}

func Refresh(c *fiber.Ctx) error {
	var request struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": false, "message": "Could not process JSON"})
	}

	validator := validator.New()

	if err := validator.Struct(request); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{"status": false, "message": "Please pass in the correct data"})
	}

	email, err := database.RedisClient.Get(request.RefreshToken)
	if err != nil {
		if err == redis.Nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"status": false, "message": "User not logged in"})
		}
	}

	user, err := services.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": false, "message": "User not found"})
		}
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": false, "message": err.Error()})
	}

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
			JSON(fiber.Map{"status": false, "message": "Could not sign access token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "access_token": accessToken})
}

func Logout(c *fiber.Ctx) error {
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": false, "message": "Error parsing JSON"})
	}

	if request.RefreshToken == "" {
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{"status": false, "message": "Please pass in the refresh token"})
	}

	if _, err := database.RedisClient.Get(request.RefreshToken); err != nil {
		if err == redis.Nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"status": false, "message": "User not logged in"})
		}
	}

	if err := database.RedisClient.Delete(request.RefreshToken); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": false, "message": "Some error occured"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "message": "Logout Successful"})
}
