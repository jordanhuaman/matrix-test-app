package models

import "gorm.io/gorm"

type User struct {
	Base
	Email    string `json:"email"    gorm:"uniqueIndex;not null"`
	Password string `json:"-"        gorm:"not null"` // "-" nunca se serializa en JSON
	Username string `json:"user_name" gorm:"not null"`
	IsActive bool   `json:"isActive" gorm:"default:true"`

	// Relación: un usuario tiene muchos resultados de matrices
	MatrixResults []MatrixResult `json:"matrixResults,omitempty" gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser adds a new user to the database
func (r *UserRepository) CreateUser(email, username, passwordHash string) (*User, error) {
	user := &User{
		Email:    email,
		Username: username,
		Password: passwordHash,
	}

	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(id string) (*User, error) {
	var user User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepository) UpdateNames(id string, names string) (*User, error) {
	var user User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	user.Username = names
	if err := r.db.Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) DeleteUser(id string) error {
	var user User
	if err := r.db.Where("id = ?", id).Delete(&user).Error; err != nil {
		return err
	}
	return nil
}
