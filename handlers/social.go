package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/padp721/padp721-web-backend/models"
)

func GetSocials(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	sql := "SELECT name, url, color, icon_type, icon FROM public.socials"
	rows, err := db.Query(c.Context(), sql)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}
	defer rows.Close()

	var socials []models.Social
	for rows.Next() {
		var social models.Social
		if err := rows.Scan(&social.Name, &social.Url, &social.Color, &social.IconType, &social.Icon); err != nil {
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

func GetSocial(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	id := c.Params("id")

	var social models.Social
	sql := "SELECT name, url, color, icon_type, icon FROM public.socials WHERE id=$1"
	err := db.QueryRow(c.Context(), sql, id).Scan(&social.Name, &social.Url, &social.Color, &social.IconType, &social.Icon)
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

func CreateSocial(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	var social models.Social
	if err := c.BodyParser(&social); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Message: err.Error(),
		})
	}

	sql := "INSERT INTO public.socials(name, url, color, icon_type, icon) VALUES($1, $2, $3, $4, $5)"
	_, err := db.Exec(c.Context(), sql, social.Name, social.Url, social.Color, social.IconType, social.Icon)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Message: "Social Created.",
	})
}

func UpdateSocial(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: "Database connection not available!",
		})
	}

	var social models.Social
	if err := c.BodyParser(&social); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Message: err.Error(),
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

func DeleteSocial(c *fiber.Ctx) error {
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
