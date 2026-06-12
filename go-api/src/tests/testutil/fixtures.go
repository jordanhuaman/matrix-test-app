package testutil

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jordanhuaman/go-api/src/models"
	"github.com/jordanhuaman/go-api/src/services"
	"gorm.io/gorm"
)

func SeedUser(t *testing.T, db *gorm.DB, email, password string) *models.User {
	t.Helper()

	hash, err := services.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	user := &models.User{
		Email:    email,
		Password: hash,
		Username: "Test User",
		IsActive: true,
	}
	user.ID = uuid.New()

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}
	return user
}

func SeedRefreshToken(t *testing.T, db *gorm.DB, userID string, ttl time.Duration) *models.RefreshToken {
	t.Helper()

	token := &models.RefreshToken{
		Token:     uuid.New().String(),
		UserId:    userID,
		ExpiresAt: time.Now().Add(ttl),
		Revoked:   false,
	}
	if err := db.Create(token).Error; err != nil {
		t.Fatalf("failed to seed refresh token: %v", err)
	}
	return token
}

func SeedMatrixResult(t *testing.T, db *gorm.DB, userID uuid.UUID) *models.MatrixResult {
	t.Helper()

	input := &models.MatrixInput{
		Data:    [][]float64{{1, 2}, {3, 4}, {5, 6}},
		Rows:    3,
		Columns: 2,
	}
	input.ID = uuid.New()
	if err := db.Create(input).Error; err != nil {
		t.Fatalf("failed to seed matrix input: %v", err)
	}

	result := &models.MatrixResult{
		UserID:        userID,
		MatrixInputID: input.ID,
		QRResult: models.QRMatrices{
			Q: [][]float64{{-0.169, 0.963}, {-0.507, 0.259}, {-0.845, -0.072}},
			R: [][]float64{{-5.916, -7.603}, {0, 1.300}},
		},
		Statistics: models.Statistics{
			Max:       7.603,
			Min:       -7.603,
			Average:   0,
			Sum:       0,
			QDiagonal: false,
			RDiagonal: false,
		},
		Status: "completed",
	}
	result.ID = uuid.New()
	if err := db.Create(result).Error; err != nil {
		t.Fatalf("failed to seed matrix result: %v", err)
	}
	return result
}

func GenerateToken(secret string, userID string) string {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}
