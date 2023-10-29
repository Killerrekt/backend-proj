package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "www.github.com/ic-ETITE-24/icetite-24-backend/internal/controllers"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/middleware"
)

func AdminRoutes(incomingRoutes *fiber.App) {
	adminRouter := incomingRoutes.Group("/admin")
	adminRouter.Use(middleware.VerifyAdminToken)

	adminRouter.Get("/users", controller.GetAllUsers)
	adminRouter.Get("/team/project/:id", controller.GetProjectFromTeamID)
	adminRouter.Get("/team/idea/:id", controller.GetIdeaFromTeamID)
	adminRouter.Get("/teams", controller.GetAllTeams)
	adminRouter.Get("/team/user/:id", controller.GetLeaderInfo)
	adminRouter.Post("/ban/:id", controller.BanUser)
	adminRouter.Post("/unban/:id", controller.UnbanUser)
}
