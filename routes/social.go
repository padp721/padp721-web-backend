package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/padp721/padp721-web-backend/handlers"
)

func SetupSocialRoutes(parentRoute fiber.Router) {
	socialRoutes := parentRoute.Group("/social")

	socialRoutes.Get("/", handlers.SocialsGet)
	socialRoutes.Get("/:id", handlers.SocialGet)
	socialRoutes.Post("/", handlers.SocialCreate)
	socialRoutes.Put("/:id", handlers.SocialUpdate)
	socialRoutes.Delete("/:id", handlers.SocialDelete)
}
