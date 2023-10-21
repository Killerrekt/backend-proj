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
	incomingRoutes.Get("/users/logout", controller.Logout)
	incomingRoutes.Post("/users/resetpassword/", controller.SendResetPasswordOTP)
	incomingRoutes.Put("/users/resetpassword/", controller.VerifyResetPasswordOTP)
	incomingRoutes.Post("/users/verifyuser/", controller.SendVerifyUserOTP)
	incomingRoutes.Put("/users/verifyuser", controller.VerifyUserOTP)
  incomingRoutes.Get("/users/getall", middleware.VerifyAdminToken, controller.GetAllUsers)

	incomingRoutes.Get("/payment/initiate", middleware.VerifyAccessToken, controller.InitiatePayment)
	incomingRoutes.Post("/payment/callbackurl", controller.CallBackURL)
}
