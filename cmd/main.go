package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"www.github.com/ic-ETITE-24/icetite-24-backend/config"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/routes"
)

func main() {
	app := fiber.New()
	config.SanityCheck()
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
	err = database.NewRepository(redisConfig)
	if err != nil {
		log.Fatalln("Failed to load redis client! \n", err.Error())
	}

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.ClientOrigin,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	app.Use(recover.New())

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Pong -> Welcome to the ICETITE 24 Hackathon Backend",
		})
	})

	routes.UserRoutes(app)
	routes.PaymentRoutes(app)
	routes.ProjectsRoutes(app)
	routes.TeamRoutes(app)
	routes.AdminRoutes(app)
	routes.EnquiryRoutes(app)
	routes.IdeasRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Route not found",
		})
	})

	log.Fatal(app.Listen(config.Port))
}
