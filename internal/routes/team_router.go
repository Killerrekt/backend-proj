package routes

import (
	"github.com/gofiber/fiber/v2"

	controllers "www.github.com/ic-ETITE-24/icetite-24-backend/internal/controllers"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/middleware"
)

func TeamRoutes(app *fiber.App) {
	teamRoutes := app.Group("/teams")
	teamRoutes.Post("/", middleware.VerifyAccessToken, controllers.CreateTeam)
	teamRoutes.Post("/join", middleware.VerifyAccessToken, controllers.JoinTeam)
	teamRoutes.Get("/", middleware.VerifyAccessToken, controllers.GetTeam)
	teamRoutes.Put("/", middleware.VerifyAccessToken, controllers.UpdateTeam)
	teamRoutes.Delete("/", middleware.VerifyAccessToken, controllers.DeleteTeam)
	teamRoutes.Get("/leave", middleware.VerifyAccessToken, controllers.LeaveTeam)
}
