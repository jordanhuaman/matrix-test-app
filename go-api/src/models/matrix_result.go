package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

// QRMatrices contiene las matrices Q y R resultado de la factorización
// Se guarda como JSONB en PostgreSQL
type QRMatrices struct {
	Q Matrix2D `json:"q"` // Matriz ortogonal
	R Matrix2D `json:"r"` // Matriz triangular superior
}

// Value implementa driver.Valuer para guardar en PostgreSQL
func (q QRMatrices) Value() (driver.Value, error) {
	return json.Marshal(q)
}

// Scan implementa sql.Scanner para leer desde PostgreSQL
func (q *QRMatrices) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("error al convertir QRMatrices desde base de datos")
	}
	return json.Unmarshal(bytes, q)
}

// Statistics contiene las estadísticas calculadas por Node.js
type Statistics struct {
	Max       float64 `json:"max"`
	Min       float64 `json:"min"`
	Average   float64 `json:"average"`
	Sum       float64 `json:"sum"`
	QDiagonal bool    `json:"qIsDiagonal"` // si Q es matriz diagonal
	RDiagonal bool    `json:"rIsDiagonal"` // si R es matriz diagonal
}

// Value implementa driver.Valuer
func (s Statistics) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan implementa sql.Scanner
func (s *Statistics) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("error al convertir Statistics desde base de datos")
	}
	return json.Unmarshal(bytes, s)
}

// MatrixResult guarda el resultado completo de una operación
type MatrixResult struct {
	Base

	// Relación con el usuario que hizo la petición
	UserID uuid.UUID `json:"userId" gorm:"type:uuid;not null;index"`
	User   User      `json:"user,omitempty" gorm:"foreignKey:UserID"`

	// Relación con la matriz de entrada
	MatrixInputID uuid.UUID   `json:"matrixInputId" gorm:"type:uuid;not null"`
	MatrixInput   MatrixInput `json:"matrixInput,omitempty" gorm:"foreignKey:MatrixInputID"`

	// Resultado de la factorización QR (calculado por Go)
	QRResult QRMatrices `json:"qrResult" gorm:"type:jsonb;not null"`

	// Estadísticas calculadas por Node.js
	Statistics Statistics `json:"statistics" gorm:"type:jsonb;not null"`

	// Estado del proceso
	Status   string `json:"status" gorm:"default:'pending'"` // pending, completed, error
	ErrorMsg string `json:"errorMsg,omitempty"`
}

func (MatrixResult) TableName() string {
	return "matrix_results"
}
