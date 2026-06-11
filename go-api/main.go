package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	database "github.com/jordanhuaman/go-api/src/framework"
)

func main() {
	godotenv.Load()
	database.ConnectDb()
	// Create new Fiber instance
	app := fiber.New()
	// Create new GET route on path "/"
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	// Start server on http://localhost:3000
	log.Fatal(app.Listen(":3000"))
}
