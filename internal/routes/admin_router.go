package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "www.github.com/ic-ETITE-24/icetite-24-backend/internal/controllers"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/middleware"
)

func AdminRoutes(incomingRoutes *fiber.App) {
	adminRouter := incomingRoutes.Group("/admin")
	adminRouter.Use(middleware.VerifyAdminToken)

	adminRouter.Get("/getall", controller.GetAllUsers)
	adminRouter.Post("/ban/:id", controller.BanUser)
	adminRouter.Post("/unban/:id", controller.UnbanUser)
}
