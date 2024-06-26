package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/padp721/padp721-web-backend/models"
)

func SocialsGet(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	sql := "SELECT id, name, url, color, icon_type, icon FROM public.socials ORDER BY created_at DESC"
	rows, err := db.Query(c.Context(), sql)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}
	defer rows.Close()

	socials := []models.Social{}
	for rows.Next() {
		var social models.Social
		if err := rows.Scan(&social.Id, &social.Name, &social.Url, &social.Color, &social.IconType, &social.Icon); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
				Message: err.Error(),
			})
		}
		socials = append(socials, social)
	}

	return c.JSON(models.ResponseData{
		Message: "Data fetch success!",
		Data:    socials,
	})
}

func SocialGet(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	id := c.Params("id")

	var social models.Social
	sql := "SELECT id, name, url, color, icon_type, icon FROM public.socials WHERE id=$1"
	err := db.QueryRow(c.Context(), sql, id).Scan(&social.Id, &social.Name, &social.Url, &social.Color, &social.IconType, &social.Icon)
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
		Data:    social,
	})
}

func SocialCreate(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	var social models.SocialInput
	if err := c.BodyParser(&social); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Message: err.Error(),
		})
	}

	validate, ok := c.Locals("validator").(*validator.Validate)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Validator not found!",
		})
	}

	if err := validate.Struct(social); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Message: fmt.Sprintf("Validation error: %v", err),
		})
	}

	sql := "INSERT INTO public.socials(id, name, url, color, icon_type, icon) VALUES($1, $2, $3, $4, $5, $6)"
	_, err := db.Exec(c.Context(), sql, uuid.New(), social.Name, social.Url, social.Color, social.IconType, social.Icon)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Message: "Social Created.",
	})
}

func SocialUpdate(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	var social models.SocialInput
	if err := c.BodyParser(&social); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Message: err.Error(),
		})
	}

	validate, ok := c.Locals("validator").(*validator.Validate)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Validator not found!",
		})
	}

	if err := validate.Struct(social); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Message: fmt.Sprintf("Validation error: %v", err),
		})
	}

	id := c.Params("id")

	sql := "UPDATE public.socials SET name=$2, url=$3, color=$4, icon_type=$5, icon=$6 WHERE id=$1"
	commandTag, err := db.Exec(c.Context(), sql, id, social.Name, social.Url, social.Color, social.IconType, social.Icon)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	if commandTag.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Message: "Social not found!",
		})
	}

	return c.JSON(models.Response{
		Message: "Social Updated.",
	})
}

func SocialDelete(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	id := c.Params("id")

	sql := "DELETE FROM public.socials WHERE id=$1"
	commandTag, err := db.Exec(c.Context(), sql, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	if commandTag.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Message: "Social not found!",
		})
	}

	return c.JSON(models.Response{
		Message: "Social Deleted.",
	})
}
