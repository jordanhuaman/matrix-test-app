package models

import "gorm.io/gorm"

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
