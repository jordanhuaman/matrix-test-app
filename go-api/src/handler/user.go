package handlers

import (
	"log"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jordanhuaman/go-api/src/models"
)

// UserHandler contains HTTP handlers for users.
type UserHandler struct {
	userRepo         *models.UserRepository
	refreshTokenRepo *models.RefreshTokenRepository
}

// NewUserHandler creates a new user handler.
func NewUserHandler(userRepo *models.UserRepository, refreshTokenRepo *models.RefreshTokenRepository) *UserHandler {
	return &UserHandler{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func validToken(t *jwt.Token, id string) bool {
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("[validToken] ERROR: failed to cast claims to MapClaims")
		return false
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		log.Printf("[validToken] ERROR: sub claim is not a string, got: %T value: %v\n", claims["sub"], claims["sub"])
		return false
	}

	match := sub == id
	log.Printf("[validToken] sub=%q | param_id=%q | match=%v\n", sub, id, match)
	return match
}

func parseUserID(c fiber.Ctx) (string, error) {
	idParam := c.Params("id")
	_, err := uuid.Parse(idParam)
	if err != nil {
		return "", c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user ID",
			"data":    nil,
		})
	}

	return idParam, nil
}

// GetUser get a user
func (uh *UserHandler) GetUser(c fiber.Ctx) error {
	id, err := parseUserID(c)
	if err != nil {
		return err
	}

	tok := jwtware.FromContext(c)
	if tok == nil {
		log.Printf("[GetUser] ERROR: token not found in context\n")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Token not found",
			"data":    nil,
		})
	}

	log.Printf("[GetUser] tok valid=%v claims=%v\n", tok.Valid, tok.Claims)

	if !validToken(tok, id) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token id",
			"data":    nil,
		})
	}

	user, err := uh.userRepo.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No user found with ID",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User found",
		"data":    user,
	})
}

// UpdateUser update user
func (uh *UserHandler) UpdateUser(c fiber.Ctx) error {
	var input struct {
		Names string `json:"names"`
	}
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
			"data":    err.Error(),
		})
	}

	id, err := parseUserID(c)
	if err != nil {
		return err
	}

	tok := jwtware.FromContext(c)
	if tok == nil || !validToken(tok, id) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token id",
			"data":    nil,
		})
	}

	updatedUser, err := uh.userRepo.UpdateNames(id, input.Names)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "User update failed",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User successfully updated",
		"data":    updatedUser,
	})
}

// DeleteUser delete user
func (uh *UserHandler) DeleteUser(c fiber.Ctx) error {
	id, err := parseUserID(c)
	if err != nil {
		return err
	}

	tok := jwtware.FromContext(c)
	if tok == nil || !validToken(tok, id) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token id",
			"data":    nil,
		})
	}

	if err := uh.refreshTokenRepo.RevokeAllUserTokens(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to revoke user tokens",
			"data":    nil,
		})
	}

	err = uh.userRepo.DeleteUser(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"data":    nil,
		})
	}

	c.Locals("user", nil)
	c.ClearCookie("jwt")

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User successfully deleted",
		"data":    nil,
	})
}
