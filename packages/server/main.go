package main

import (
	routes "area-server/api"
	config "area-server/config"
	"area-server/db/postgres"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// `LoadApplets()` loads all the triggers from the database and stores them in the `applets` map
func main() {
	// Database
	pg := postgres.Init()
	pg.Connect()

	if !pg.OK() {
		// Auto migrate if not exist
		pg.Migrate()
	}

	if _, ok := os.LookupEnv("AREA_STATE"); !ok {
		panic("AREA_STATE not set")
	}

	// Load triggers - If exist
	if err := LoadApplets(); err != nil {
		panic(err)
	}

	// Create folder avatars if not exist
	os.Mkdir("./avatars", 666)

	// Fiber
	app := fiber.New()

	// General Middlewares
	app.Use(logger.New())
	if config.CFG.Mode == config.Session {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     "http://localhost:8081",
			AllowCredentials: true,
			AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Session",
		}))
	} else {
		app.Use(cors.New())
	}

	// API
	routes.Init(app)

	if config.CFG.HTTPS {
		app.ListenTLS(":8080", "./ssl/cert.cert", "./ssl/cert.key")
	}
	app.Listen(":8080")
}
