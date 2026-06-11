package models

import "time"

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
