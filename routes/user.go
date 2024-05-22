package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/padp721/padp721-web-backend/handlers"
)

func SetupUserRoutes(parentRoute fiber.Router) {
	userRoutes := parentRoute.Group("/user")

	// userRoutes.Get("/")
	// userRoutes.Get("/:id")
	userRoutes.Post("/", handlers.UserCreate)
	// userRoutes.Put("/:id")
	// userRoutes.Delete("/:id")
}
