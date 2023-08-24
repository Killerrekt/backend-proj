package routes

import (
	"github.com/gofiber/fiber/v2"

	controller "www.github.com/ic-ETITE-24/icetite-24-backend/controllers"
)

func UserRoutes(incomingRoutes *fiber.App) {
	incomingRoutes.Post("/users/signup", controller.CreateUser)
}
