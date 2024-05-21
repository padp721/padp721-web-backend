package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/padp721/padp721-web-backend/handlers"
)

func SetupSocialRoutes(api fiber.Router) {
	socialRoutes := api.Group("/social")

	socialRoutes.Get("/", handlers.GetSocials)
	socialRoutes.Post("/", handlers.CreateSocial)
}
