package main

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/jordanhuaman/go-api/src/database"
	"github.com/jordanhuaman/go-api/src/router"
)

// 1. Creamos una estructura propia que envuelve al validador original
type GoPlaygroundValidator struct {
	Validator *validator.Validate
}

// 2. Implementamos el método Validate que exige la interfaz fiber.StructValidator de v3
func (v *GoPlaygroundValidator) Validate(out any) error {
	return v.Validator.Struct(out)
}
func main() {
	// godotenv.Load()

	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Warning: no .env file found, using environment variables")
	}
	database.ConnectDb()

	cv := &GoPlaygroundValidator{Validator: validator.New()}

	app := fiber.New(fiber.Config{StructValidator: cv})
	app.Use(cors.New(cors.Config{AllowOrigins: []string{"http://localhost:3000"}}))
	router.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
