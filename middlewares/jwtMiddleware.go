package middlewares

import (
	"crypto/rsa"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/padp721/padp721-web-backend/models"
	"github.com/padp721/padp721-web-backend/utilities"
)

func JwtMiddleware(c *fiber.Ctx) error {
	jwtTokenString := c.Cookies("jwt")
	if jwtTokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Message: "JWT Token not found!",
		})
	}

	privateKey, ok := c.Locals("privateKey").(*rsa.PrivateKey)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Private key for Signing JWT is not loaded.",
		})
	}

	token, err := utilities.ParseJWT(privateKey, jwtTokenString)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: fmt.Sprintf("Cannot parse JWT Key: %v", err),
		})
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: fmt.Sprintf("Cannot get subject from JWT Key: %v", err),
		})
	}

	c.Locals("username", subject)

	return c.Next()
}
