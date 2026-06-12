package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/jordanhuaman/go-api/src/framework"
	database "github.com/jordanhuaman/go-api/src/framework"
)

// 1. Creamos una estructura propia que envuelve al validador original
type GoPlaygroundValidator struct {
	Validator *validator.Validate
}

// 2. Implementamos el método Validate que exige la interfaz fiber.StructValidator de v3
func (v *GoPlaygroundValidator) Validate(out any) error {
	return v.Validator.Struct(out)
}

func welcome(c fiber.Ctx) error {
	return c.SendString("Welcome to my API")
}

func setupRoutes(app *fiber.App) {
	app.Get("/api", welcome)
	app.Post("/user/register", framework.CreateUser)
}

func main() {
	godotenv.Load()
	database.ConnectDb()

	cv := &GoPlaygroundValidator{Validator: validator.New()}

	app := fiber.New(fiber.Config{StructValidator: cv})
	setupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
