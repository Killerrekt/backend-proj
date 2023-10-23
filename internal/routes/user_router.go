package routes

import (
	"github.com/gofiber/fiber/v2"

	controller "www.github.com/ic-ETITE-24/icetite-24-backend/internal/controllers"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/middleware"
)

func UserRoutes(incomingRoutes *fiber.App) {
	userRoutes := incomingRoutes.Group("/users")
	userRoutes.Post("/refresh", middleware.VerifyRefreshToken, controller.Refresh)
	userRoutes.Post("/signup", controller.CreateUser)
	userRoutes.Post("/login", controller.Login)
	userRoutes.Get("/logout", controller.Logout)
	userRoutes.Post("/resetpassword/", controller.SendResetPasswordOTP)
	userRoutes.Put("/resetpassword/", controller.VerifyResetPasswordOTP)
	userRoutes.Post("/verifyuser/", controller.SendVerifyUserOTP)
	userRoutes.Put("/verifyuser", controller.VerifyUserOTP)
	userRoutes.Get("/getall", middleware.VerifyAdminToken, controller.GetAllUsers)
	userRoutes.Post("/roast", middleware.VerifyAdminToken, controller.RoastUser)
	userRoutes.Post("/revoke_roast", middleware.VerifyAdminToken, controller.RevokeRoast)
}
