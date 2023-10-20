package routes

import (
	"github.com/gofiber/fiber/v2"

	controller "www.github.com/ic-ETITE-24/icetite-24-backend/internal/controllers"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/middleware"
)

func PaymentRoutes(incomingRoutes *fiber.App) {
	incomingRoutes.Get("/payment/initiate", middleware.VerifyAccessToken, controller.InitiatePayment)
	incomingRoutes.Post("/payment/callbackurl", controller.CallBackURL)
}
