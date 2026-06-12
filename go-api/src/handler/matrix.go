package handlers

import (
	"errors"
	"strconv"

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

func extractUserID(c fiber.Ctx) (uuid.UUID, error) {
	tok := jwtware.FromContext(c)
	if tok == nil {
		return uuid.Nil, fiber.ErrUnauthorized
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fiber.ErrUnauthorized
	}

	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, fiber.ErrUnauthorized
	}

	return uuid.Parse(userIDStr)
}

func (h *MatrixHandler) ListResults(c fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error", "message": "Unauthorized", "data": nil,
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	results, total, err := h.matrixService.GetUserResults(c.Context(), userID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": err.Error(), "data": nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Results retrieved",
		"data": fiber.Map{
			"results": results,
			"total":   total,
			"page":    page,
			"limit":   limit,
		},
	})
}

func (h *MatrixHandler) GetResult(c fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error", "message": "Unauthorized", "data": nil,
		})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Invalid result ID", "data": nil,
		})
	}

	result, err := h.matrixService.GetResultByID(c.Context(), userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error", "message": "Result not found", "data": nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Result found",
		"data":    result,
	})
}

func (h *MatrixHandler) CreateQR(c fiber.Ctx) error {
	var input QRRequest
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Invalid request body", "data": nil,
		})
	}

	userID, err := extractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error", "message": "Unauthorized", "data": nil,
		})
	}

	result, err := h.matrixService.ProcessMatrix(c.Context(), userID, input.Data)
	if err != nil {
		var validationErr *services.ValidationError
		if errors.As(err, &validationErr) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error", "message": err.Error(), "data": nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": err.Error(), "data": nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "QR factorization completed",
		"data":    result,
	})
}
