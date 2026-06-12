package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MatrixInputRepository struct {
	db *gorm.DB
}

func NewMatrixInputRepository(db *gorm.DB) *MatrixInputRepository {
	return &MatrixInputRepository{db: db}
}

func (r *MatrixInputRepository) Create(input *MatrixInput) error {
	return r.db.Create(input).Error
}

type MatrixResultRepository struct {
	db *gorm.DB
}

func NewMatrixResultRepository(db *gorm.DB) *MatrixResultRepository {
	return &MatrixResultRepository{db: db}
}

func (r *MatrixResultRepository) Create(result *MatrixResult) error {
	return r.db.Create(result).Error
}

func (r *MatrixResultRepository) FindByUserID(userID uuid.UUID, page, limit int) ([]MatrixResult, int64, error) {
	var results []MatrixResult
	var total int64

	query := r.db.Model(&MatrixResult{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("MatrixInput").Order("created_at DESC").Offset(offset).Limit(limit).Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *MatrixResultRepository) FindByID(userID, id uuid.UUID) (*MatrixResult, error) {
	var result MatrixResult
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).Preload("MatrixInput").First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
