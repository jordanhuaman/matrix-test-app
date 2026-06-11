package models

type User struct {
	Base
	Email    string `json:"email"    gorm:"uniqueIndex;not null"`
	Password string `json:"-"        gorm:"not null"` // "-" nunca se serializa en JSON
	FullName string `json:"fullName" gorm:"not null"`
	IsActive bool   `json:"isActive" gorm:"default:true"`

	// Relación: un usuario tiene muchos resultados de matrices
	MatrixResults []MatrixResult `json:"matrixResults,omitempty" gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}
