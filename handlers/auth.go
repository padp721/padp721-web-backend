package handlers

import (
	"crypto/rsa"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/padp721/padp721-web-backend/models"
	"github.com/padp721/padp721-web-backend/utilities"
	"golang.org/x/crypto/bcrypt"
)

func AuthLogin(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	var auth models.Auth
	if err := c.BodyParser(&auth); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Message: err.Error(),
		})
	}

	var user models.UserSecret
	sql := "SELECT id, username, password, name FROM public.users WHERE username=$1"
	err := db.QueryRow(c.Context(), sql, auth.Username).Scan(&user.Id, &user.Username, &user.Password, &user.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
				Message: "User tidak ditemukan!",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	passwordString := utilities.GeneratePasswordString(
		auth.Password,
		user.Username,
		user.Name,
	)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordString))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Message: "Password Salah!",
		})
	}

	privateKey, ok := c.Locals("privateKey").(*rsa.PrivateKey)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Private key for Signing JWT is not loaded.",
		})
	}
	jwtTokenString, jwtExpire, err := utilities.GenerateJWT(privateKey, user.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: fmt.Sprintf("Failed to create JWT Token: %v", err),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:    "jwt",
		Value:   jwtTokenString,
		Expires: jwtExpire,
	})

	return c.JSON(models.Response{
		Message: "Berhasil Login.",
	})
}
