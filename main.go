package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/padp721/padp721-web-backend/db"
	"github.com/padp721/padp721-web-backend/routes"
)

var devMode bool

// * LOAD ENV VARIABLES
func init() {
	log.Println("Initializing App...")
	if devMode = strings.ToLower(os.Getenv("APP_DEV_MODE")) != "false"; devMode {
		log.Println("App Running in Development Mode.")
		log.Println("Loading from .env file...")
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		log.Println("App Running in Production Mode.")
	}
}

func main() {
	App := fiber.New(fiber.Config{
		AppName: os.Getenv("APP_NAME"),
	})
	DbPool := db.Connect()
	defer DbPool.Close()

	if devMode {
		App.Use(logger.New())
	}
	App.Use(cors.New())
	App.Use(func(c *fiber.Ctx) error {
		c.Locals("db", DbPool)
		return c.Next()
	})

	App.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"Hello": "World!"})
	})

	api := App.Group("/api")
	routes.SetupSocialRoutes(api)

	APP_HOST := os.Getenv("APP_HOST")
	APP_PORT := os.Getenv("APP_PORT")
	log.Fatal(App.Listen(fmt.Sprintf("%v:%v", APP_HOST, APP_PORT)))
}
