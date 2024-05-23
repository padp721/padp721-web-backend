package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/padp721/padp721-web-backend/handlers"
)

func SetupAuthRoutes(app *fiber.App) {
	authRoutes := app.Group("/auth")

	authRoutes.Post("/login", handlers.AuthLogin)
	authRoutes.Get("/logout", handlers.AuthLogout)
}
