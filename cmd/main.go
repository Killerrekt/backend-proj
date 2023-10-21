package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"www.github.com/ic-ETITE-24/icetite-24-backend/config"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/routes"
)

func main() {
	app := fiber.New()

	redisConfig, err := config.LoadRedisConfig()

	if err != nil {
		log.Fatalln("Failed to load redis environment variable! \n", err.Error())
	}

	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}

	database.ConnectDB(&config)
	database.RunMigrations(database.DB)
	fmt.Println(redisConfig)
	err = database.NewRepository(redisConfig)
	if err != nil {
		log.Fatalln("Failed to load redis client! \n", err.Error())
	}

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.ClientOrigin,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	apiGroup := app.Group("/v1")

	// routes.AuthRoutes(apiGroup)

	apiGroup.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "icETITE-24 Backend API is up and running.",
		})
	})

	routes.UserRoutes(app)
	routes.ProjectsRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Route not found",
		})
	})

	log.Fatal(app.Listen(config.Port))
}
