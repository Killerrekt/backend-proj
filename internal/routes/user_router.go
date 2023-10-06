package routes

import (
	"github.com/gofiber/fiber/v2"

	controller "www.github.com/ic-ETITE-24/icetite-24-backend/internal/controllers"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/middleware"
)

func UserRoutes(incomingRoutes *fiber.App) {
	incomingRoutes.Post("/users/refresh", middleware.VerifyRefreshToken, controller.Refresh)
	incomingRoutes.Post("/users/signup", controller.CreateUser)
	incomingRoutes.Post("/users/login", controller.Login)
	incomingRoutes.Get("/users/logout", middleware.VerifyAccessToken, controller.Logout)
	incomingRoutes.Get("/users/forgotpassword/:email", controller.ForgotPassword)
	incomingRoutes.Post("/users/resetpassword/:Token", controller.ResetPassword)
}
