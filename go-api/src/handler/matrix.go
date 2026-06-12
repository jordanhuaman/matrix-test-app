package handlers

import (
	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jordanhuaman/go-api/src/services"
)

type MatrixHandler struct {
	matrixService *services.MatrixService
}

func NewMatrixHandler(matrixService *services.MatrixService) *MatrixHandler {
	return &MatrixHandler{matrixService: matrixService}
}

type QRRequest struct {
	Data [][]float64 `json:"data" validate:"required"`
}

func (h *MatrixHandler) CreateQR(c fiber.Ctx) error {
	var input QRRequest
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"data":    nil,
		})
	}

	tok := jwtware.FromContext(c)
	if tok == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
			"data":    nil,
		})
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token claims",
			"data":    nil,
		})
	}

	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token subject",
			"data":    nil,
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user ID in token",
			"data":    nil,
		})
	}

	result, err := h.matrixService.ProcessMatrix(c.Context(), userID, input.Data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "QR factorization completed",
		"data":    result,
	})
}
