package middleware

import (
	"errors"
	"log"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/services"
)

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(viper.GetString("ACCESS_SECRET_KEY"))},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Println(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "Invalid or expired JWT"})
		},
	})
}

func VerifyAccessToken(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	person, err := services.FindUserByEmail(claims["sub"].(string))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": false, "message": "User not found"})
		}
	}

	if person.IsBanned {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "User is banned"})
	}

	c.Locals("user", person)
	return c.Next()
}

func VerifyAdminToken(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	if claims["role"] != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false, "message": "Invalid Role",
		})
	}

	var user models.User
	database.DB.Find(&user, "email = ?", claims["sub"])

	if user.ID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false, "message": "Invalid User",
		})
	}

	if user.IsBanned {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false, "message": "User is banned", "banned": true,
		})
	}

	c.Locals("user", user)
	return c.Next()
}
