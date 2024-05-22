package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/padp721/padp721-web-backend/models"
	"github.com/padp721/padp721-web-backend/utilities"
	"golang.org/x/crypto/bcrypt"
)

func UserCreate(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	var user models.UserInput
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Message: err.Error(),
		})
	}

	passwordString := utilities.GeneratePasswordString(
		user.Password,
		user.Username,
		user.Name,
	)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordString), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	sql := "INSERT INTO public.users(username, password, name) VALUES($1, $2, $3)"
	_, err = db.Exec(c.Context(), sql, user.Username, string(hashedPassword), user.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Message: "User Created.",
	})
}

func UsersGet(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	sql := "SELECT id, name, username FROM public.users"
	rows, err := db.Query(c.Context(), sql)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Username); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
				Message: err.Error(),
			})
		}
		users = append(users, user)
	}

	return c.JSON(models.ResponseData{
		Message: "Data fetch success!",
		Data:    users,
	})
}

func UserGet(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	id := c.Params("id")

	var user models.User
	sql := "SELECT id, name, username FROM public.users WHERE id=$1"
	err := db.QueryRow(c.Context(), sql, id).Scan(&user.Id, &user.Name, &user.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(models.Response{
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(models.ResponseData{
		Message: "Data fetch success!",
		Data:    user,
	})
}
