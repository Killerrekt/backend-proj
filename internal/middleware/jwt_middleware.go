package middleware

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"

	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

func VerifyAccessToken(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"status": false, "message": "Missing Authorization Header"})
	}

	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"status": false, "message": "Invalid Authorization Header Format"})
	}

	token := strings.TrimPrefix(tokenString, "Bearer ")

	accessToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(viper.GetString("ACCESS_SECRET_KEY")), nil
	})

	if claims, ok := accessToken.Claims.(jwt.MapClaims); ok && accessToken.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": false, "message": "Token Expired",
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
				"status": false, "message": "User is banned", "roasted": true,
			})
		}

		c.Locals("user", user)
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "Invalid Token"})
}

func VerifyAdminToken(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"status": false, "message": "Missing Authorization Header"})
	}

	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"status": false, "message": "Invalid Authorization Header Format"})
	}

	token := strings.TrimPrefix(tokenString, "Bearer ")

	accessToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(viper.GetString("ACCESS_SECRET_KEY")), nil
	})

	if claims, ok := accessToken.Claims.(jwt.MapClaims); ok && accessToken.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": false, "message": "Token Expired",
			})
		}

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
				"status": false, "message": "User is banned", "roasted": true,
			})
		}

		c.Locals("user", user)
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "Invalid Token"})
}

func VerifyRefreshToken(c *fiber.Ctx) error {
	var tokenReq struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&tokenReq); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": false, "message": "Please pass in the refresh_token"})
	}

	if tokenReq.RefreshToken == "" {
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{"status": false, "message": "Please pass in the refresh token"})
	}

	token := tokenReq.RefreshToken

	email, err := database.RedisClient.Get(tokenReq.RefreshToken)

	if err == redis.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "User not logged in",
		})
	}

	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": false, "message": "internal server error"})
	}

	accessToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(viper.GetString("REFRESH_SECRET_KEY")), nil
	})

	if claims, ok := accessToken.Claims.(jwt.MapClaims); ok && accessToken.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": false, "message": "Token Expired",
			})
		}

		var user models.User
		database.DB.Find(&user, "email = ?", email)

		if user.ID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": false, "message": "Invalid User",
			})
		}

		c.Locals("user", user)
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "Invalid Token"})
}
