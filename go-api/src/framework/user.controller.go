package framework

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jordanhuaman/go-api/src/models"
)

type CreateUserDTO struct {
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"fullName" validate:"required"`
	Password string `json:"password" validate:"required,min=6"` // Obligatorio y mínimo 6 caracteres
}

// func CreateUserResponse(user models.User) User {
// 	return User{ID: user.ID.String(), Email: user.Email, FullName: user.FullName}
// }

func CreateUser(c fiber.Ctx) error {
	var dto CreateUserDTO

	if err := c.Bind().Body(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Datos e formato inválido: " + err.Error(),
		})
	}

	user := models.User{
		Email:    dto.Email,
		FullName: dto.FullName,
		Password: dto.Password,
	}

	result := Database.Db.Create(&user)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "No se pudo crear el usuario: " + result.Error.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"id":      user.ID,
		"message": "Usuario creado con éxito",
	})
}
