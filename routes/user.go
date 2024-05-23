package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/padp721/padp721-web-backend/handlers"
)

func SetupUserRoutes(parentRoute fiber.Router) {
	userRoutes := parentRoute.Group("/user")

	userRoutes.Get("/", handlers.UsersGet)
	userRoutes.Get("/:id", handlers.UserGet)
	userRoutes.Post("/", handlers.UserCreate)
	userRoutes.Put("/:id", handlers.UserUpdate)
	userRoutes.Delete("/:id", handlers.UserDelete)
	userRoutes.Patch("/change-password", handlers.UserChangePassword)
}
