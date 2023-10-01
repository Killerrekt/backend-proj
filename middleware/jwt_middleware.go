package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"

	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

func VerifyAccessToken(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "Missing Authorization Header"})
	}

	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "Invalid Authorization Header Format"})
	}

	token := strings.TrimPrefix(tokenString, "Bearer ")

	accessToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signing method")
		}
		return []byte(viper.GetString("ACCESS_SECRET_KEY")), nil
	})

	if claims, ok := accessToken.Claims.(jwt.MapClaims); ok && accessToken.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token Expired"})
		}

		var user models.User
		database.DB.Find(&user, "email = ?", claims["sub"])

		if user.ID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid User"})
		}

		if float64(user.TokenVersion) != claims["version"].(float64) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Token"})
		}

		c.Locals("user", user)
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Token"})
}

func VerifyAdminToken(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "Missing Authorization Header"})
	}

	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "Invalid Authorization Header Format"})
	}

	token := strings.TrimPrefix(tokenString, "Bearer ")

	accessToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signing method")
		}
		return []byte(viper.GetString("ACCESS_SECRET_KEY")), nil
	})

	if claims, ok := accessToken.Claims.(jwt.MapClaims); ok && accessToken.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token Expired"})
		}

		if claims["role"] != "admin" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Role"})
		}

		var user models.User
		database.DB.Find(&user, "email = ?", claims["sub"])

		if user.ID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid User"})
		}

		if user.TokenVersion != claims["version"].(int) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Token"})
		}

		c.Locals("user", user)
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Token"})
}

func VerifyRefreshToken(c *fiber.Ctx) error {
	var tokenReq struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&tokenReq); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "Please pass in the refresh_token"})
	}

	token := tokenReq.RefreshToken

	accessToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signing method")
		}
		return []byte(viper.GetString("REFRESH_SECRET_KEY")), nil
	})

	if claims, ok := accessToken.Claims.(jwt.MapClaims); ok && accessToken.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token Expired"})
		}

		var user models.User
		database.DB.Find(&user, "email = ?", claims["sub"])

		if user.ID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid User"})
		}

		if user.TokenVersion == 0 {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"message": "User not logged in"})
		}

		c.Locals("user", user)
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Token"})
}
