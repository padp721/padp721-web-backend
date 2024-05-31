package handlers

import (
	"crypto/rsa"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
			Message: fmt.Sprintf("Error parsing request body: %v", err),
		})
	}

	validate, ok := c.Locals("validator").(*validator.Validate)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Validator not found!",
		})
	}

	if err := validate.Struct(auth); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Message: fmt.Sprintf("Validation error: %v", err),
		})
	}

	//* FIND USER
	var (
		userId     uuid.UUID
		username   string
		dbPassword string
	)
	sql := `
		SELECT a.id, a.username, b.password 
		FROM public.users AS a
		INNER JOIN secret.users AS b ON a.id = b.user_id
		WHERE username=$1
	`
	err := db.QueryRow(c.Context(), sql, auth.Username).Scan(&userId, &username, &dbPassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
				Message: "User tidak ditemukan!",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: fmt.Sprintf("Error getting user data from database: %v", err),
		})
	}

	//* VERIFY PASSWORD
	inputPassword := utilities.GeneratePasswordString(auth.Password, username)
	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(inputPassword))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Message: "Wrong Password!",
		})
	}

	//* GENERATE JWT TOKEN
	privateKey, ok := c.Locals("privateKey").(*rsa.PrivateKey)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Private key for Signing JWT is not loaded.",
		})
	}
	jwtTokenString, _, err := utilities.GenerateJWT(privateKey, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: fmt.Sprintf("Failed to create JWT Token: %v", err),
		})
	}

	return c.JSON(models.ResponseData{
		Message: "Login Success.",
		Data: fiber.Map{
			"token": jwtTokenString,
		},
	})
}
