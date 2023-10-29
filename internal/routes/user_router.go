package routes

import (
	"github.com/gofiber/fiber/v2"

	controller "www.github.com/ic-ETITE-24/icetite-24-backend/internal/controllers"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/middleware"
)

func UserRoutes(incomingRoutes *fiber.App) {
	userRouter := incomingRoutes.Group("/users")
	userRouter.Post("/refresh", middleware.VerifyRefreshToken, controller.Refresh)
	userRouter.Post("/signup", controller.CreateUser)
	userRouter.Post("/login", controller.Login)
	userRouter.Get("/logout", controller.Logout)
	userRouter.Post("/forgot", controller.SendForgotPasswordOTP)
	userRouter.Patch("/forgot", controller.VerifyForgotPasswordOTP)
	userRouter.Post("/verify", controller.SendVerifyUserOTP)
	userRouter.Patch("/verify", controller.VerifyUserOTP)
	userRouter.Get("/me", middleware.VerifyAccessToken, controller.UserDashboard)
	userRouter.Post("/reset-pass", middleware.VerifyAccessToken, controller.ResetPassword)
	userRouter.Patch("/update", middleware.VerifyAccessToken, controller.UpdateUser)
	userRouter.Delete("/delete", middleware.VerifyAccessToken, controller.DeleteUser)
}
