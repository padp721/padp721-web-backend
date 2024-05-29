package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/padp721/padp721-web-backend/db"
	"github.com/padp721/padp721-web-backend/handlers"
	"github.com/padp721/padp721-web-backend/middlewares"
	"github.com/padp721/padp721-web-backend/routes"
	"github.com/padp721/padp721-web-backend/utilities"
)

var (
	err        error
	devMode    bool
	privateKey *rsa.PrivateKey
	dbPool     *pgxpool.Pool
	validate   *validator.Validate
)

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

// * CREATE DB POOL
func init() {
	dbPool, err = db.CreatePool()
	if err != nil {
		log.Fatalf("Failed to Create database Pool: %v", err)
	}
}

// * LOAD PRIVATE KEY FOR JWT
func init() {
	privateKey, err = utilities.LoadPrivateKey("./keys/id_rsa")
	if err != nil {
		log.Fatalf("Failed to load RPivate Key file: %v", err)
	}
}

// * INIT VALIDATOR
func init() {
	validate = validator.New()
}

func main() {
	defer dbPool.Close()
	App := fiber.New(fiber.Config{
		AppName: os.Getenv("APP_NAME"),
	})

	if devMode {
		App.Use(logger.New())
	}
	App.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CORS_ALLOWED_ORIGIN"),
	}))
	App.Use(func(c *fiber.Ctx) error {
		c.Locals("db", dbPool)
		c.Locals("privateKey", privateKey)
		c.Locals("validator", validate)
		return c.Next()
	})

	App.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"Hello": "World!"})
	})

	keyAuthConfig := keyauth.Config{
		Validator:    middlewares.JwtMiddleware,
		ErrorHandler: middlewares.JwtOnFailure,
	}

	//* SETUP AUTH ROUTES
	routes.SetupAuthRoutes(App)

	//* SETUP BACKOFFICE ROUTES
	backOffice := App.Group("/b")
	backOffice.Use(keyauth.New(keyAuthConfig))
	routes.SetupUserRoutes(backOffice)
	routes.SetupSocialRoutes(backOffice)

	//* SETUP FRONTOFFICE ROUTES
	frontOffice := App.Group("/f")
	frontOffice.Get("/socials", handlers.SocialsGet)

	//* SERVE APP
	log.Fatal(App.Listen(fmt.Sprintf("%v:%v", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))))
}
