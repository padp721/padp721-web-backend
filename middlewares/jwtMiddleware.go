package middlewares

import (
	"crypto/rsa"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/padp721/padp721-web-backend/models"
	"github.com/padp721/padp721-web-backend/utilities"
)

func JwtOnFailure(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
		Message: err.Error(),
	})
}

func JwtMiddleware(c *fiber.Ctx, tokenString string) (bool, error) {
	if tokenString == "" {
		return false, fmt.Errorf("JWT Token not found!")
	}

	privateKey, ok := c.Locals("privateKey").(*rsa.PrivateKey)
	if !ok {
		return false, fmt.Errorf("Private key for Signing JWT is not loaded.")
	}

	token, err := utilities.ParseJWT(privateKey, tokenString)
	if err != nil {
		return false, fmt.Errorf("Cannot parse JWT Key: %v", err)
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return false, fmt.Errorf("Cannot get subject from JWT Key: %v", err)
	}

	c.Locals("userId", subject)

	return true, nil
}
