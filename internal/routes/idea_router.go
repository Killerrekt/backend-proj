package routes

import (
	"github.com/gofiber/fiber/v2"

	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/controllers"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/middleware"
)

func IdeasRoutes(incomingRoutes *fiber.App) {
	ideaRoutes := incomingRoutes.Group("/idea")
	ideaRoutes.Post("/create", middleware.VerifyAccessToken, controllers.CreateIdea)
}
