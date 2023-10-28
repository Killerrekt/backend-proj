package routes

import (
	"github.com/gofiber/fiber/v2"

	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/controllers"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/middleware"
)

func ProjectsRoutes(incomingRoutes *fiber.App) {
	projectRoutes := incomingRoutes.Group("/project")
	projectRoutes.Post("/testing", middleware.VerifyAccessToken, controllers.CreateTeam)
	projectRoutes.Post("/finalise", middleware.VerifyAccessToken, controllers.FinaliseProject)
	projectRoutes.Post("/create", middleware.VerifyAccessToken, controllers.CreateProject)
	projectRoutes.Get("/get", middleware.VerifyAccessToken, controllers.GetProject)
	projectRoutes.Delete("/delete", middleware.VerifyAccessToken, controllers.DeleteProject)
<<<<<<< HEAD
	// projectRoutes.Get("/getall", middleware.VerifyAccessToken, controllers.GetAllProject)
=======
	projectRoutes.Get("/getall", middleware.VerifyAdminToken, controllers.GetAllProject)
	projectRoutes.Post("/update", middleware.VerifyAccessToken, controllers.UpdateProject)
>>>>>>> 919f34b4797523e6af252a69b4a48588fd5be578
}
