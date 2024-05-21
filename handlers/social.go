package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/padp721/padp721-web-backend/models"
)

func GetSocials(c *fiber.Ctx) error {
	data := []models.Social{
		{
			Name:     "LinkedIn",
			Url:      "https://www.linkedin.com/in/padp721/",
			Color:    "#007bb5",
			IconType: "fab",
			Icon:     "linkedin",
		},
		{
			Name:     "Github",
			Url:      "https://github.com/padp721",
			Color:    "#6e5494",
			IconType: "fab",
			Icon:     "github",
		},
		{
			Name:     "Instagram",
			Url:      "https://www.instagram.com/padp721/",
			Color:    "null",
			IconType: "fab",
			Icon:     "instagram",
		},
		{
			Name:     "Facebook",
			Url:      "https://www.facebook.com/padp721",
			Color:    "#3B5998",
			IconType: "fab",
			Icon:     "facebook",
		},
		{
			Name:     "X",
			Url:      "https://twitter.com/padp721",
			Color:    "black",
			IconType: "fab",
			Icon:     "x-twitter",
		},
		{
			Name:     "Threads",
			Url:      "https://www.threads.net/@padp721",
			Color:    "#101010",
			IconType: "fab",
			Icon:     "threads",
		},
		{
			Name:     "E-Mail",
			Url:      "mailto:anggadp91@hotmail.com",
			Color:    "crimson",
			IconType: "fas",
			Icon:     "envelope",
		},
		{
			Name:     "Telegram",
			Url:      "https://t.me/padp721",
			Color:    "#24a6ea",
			IconType: "fab",
			Icon:     "telegram",
		},
		{
			Name:     "Discord",
			Url:      "https://discordapp.com/users/318630119373799426",
			Color:    "#545df4",
			IconType: "fab",
			Icon:     "discord",
		},
		{
			Name:     "Spotify",
			Url:      "https://open.spotify.com/user/randomize721?si=8ebc39a59bbe4243",
			Color:    "#2ada5d",
			IconType: "fab",
			Icon:     "spotify",
		},
		{
			Name:     "Steam",
			Url:      "https://steamcommunity.com/id/randomize721/",
			Color:    "#182a43",
			IconType: "fab",
			Icon:     "steam",
		},
	}

	return c.Status(200).JSON(data)
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
	_, err := db.Exec(context.Background(), sql, social.Name, social.Url, social.Color, social.IconType, social.Icon)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Message: "Social Created.",
	})
}
