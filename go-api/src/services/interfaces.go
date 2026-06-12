package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jordanhuaman/go-api/src/models"
)

type UserRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	CreateUser(email, username, passwordHash string) (*models.User, error)
}

type RefreshTokenRepository interface {
	CreateRefreshToken(userId string, ttl time.Duration) (*models.RefreshToken, error)
	GetRefreshToken(tokenString string) (*models.RefreshToken, error)
	RevokeRefreshToken(tokenString string) error
	RevokeAllUserTokens(userId string) error
}

type MatrixInputRepository interface {
	Create(input *models.MatrixInput) error
}

type MatrixResultRepository interface {
	Create(result *models.MatrixResult) error
	FindByUserID(userID uuid.UUID, page, limit int) ([]models.MatrixResult, int64, error)
	FindByID(userID, id uuid.UUID) (*models.MatrixResult, error)
}

type NodeClient interface {
	CalculateStatistics(ctx context.Context, q, r models.Matrix2D) (*models.Statistics, error)
}
