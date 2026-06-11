package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Matrix2D [][]float64

// Value convierte Matrix2D a JSON para guardar en PostgreSQL
func (m Matrix2D) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan convierte JSON de PostgreSQL a Matrix2D
func (m *Matrix2D) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("error al convertir matrix desde base de datos")
	}
	return json.Unmarshal(bytes, m)
}

// MatrixInput representa la matriz de entrada enviada por el usuario
type MatrixInput struct {
	Base

	// La matriz original enviada por el usuario
	Data Matrix2D `json:"data" gorm:"type:jsonb;not null"`

	// Dimensiones para consultas rápidas sin deserializar el JSON
	Rows    int `json:"rows"    gorm:"not null"`
	Columns int `json:"columns" gorm:"not null"`

	// Relación con el resultado generado
	MatrixResult *MatrixResult `json:"result,omitempty" gorm:"foreignKey:MatrixInputID"`
}

func (MatrixInput) TableName() string {
	return "matrix_inputs"
}
