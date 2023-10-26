package routes

import (
	"github.com/gofiber/fiber/v2"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/controllers"
)

func EnquiryRoutes(incomingRoutes *fiber.App) {
	incomingRoutes.Post("/enquiry", controllers.ExhibitionEnquiry)
}
