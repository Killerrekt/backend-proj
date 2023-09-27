package controller

import (
	"log"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid Request"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
		"sub":  user.Email,
		"role": user.Role,
	})

	accessToken, err := token.SignedString([]byte(viper.GetString("ACCESS_SECRET_KEY")))
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "Could not sign access token"})
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"sub":  user.Email,
		"role": user.Role,
	})
	refreshToken, err := token.SignedString([]byte(viper.GetString("REFRESH_SECRET_KEY")))
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "Could not sign refresh token"})
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "Login Successful", "Access Token": accessToken, "Refresh Token": refreshToken})
}
