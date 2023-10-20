package routes

import (
	"github.com/gofiber/fiber/v2"

	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/controllers"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/middleware"
)

func ProjectsRoutes(incomingRoutes *fiber.App) {
	incomingRoutes.Post("/project/create", middleware.VerifyAccessToken, controllers.CreateProject)
}
