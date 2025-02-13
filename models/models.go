package models

import (
	"github.com/lib/pq"
)

type User struct {
	ID          uint           `gorm:"primaryKey"`
	FullName    string         `json:"full_name"`
	Email       string         `json:"email"`
	Password    []byte         `json:"-"` // Store password as []byte (to avoid storing it in the response)
	PhoneNumber string         `json:"phone_number"`
	DateOfBirth string         `json:"date_of_birth"`
	City        string         `json:"city"`
	Country     string         `json:"country"`
	Role        string         `json:"role"`
	Skills      pq.StringArray `gorm:"type:text[]" json:"skills"`
}

// type User struct {
// 	Role        string   `json:"role"`
// 	FullName    string   `json:"full_name"`
// 	Email       string   `json:"email"`
// 	Password    string   `json:"password"`
// 	Skills      []string `json:"skills"`
// 	Country     string   `json:"country"`
// 	City        string   `json:"city"`
// 	DateOfBirth string   `json:"date_of_birth"`
// 	PhoneNumber string   `json:"phone_number"`
// }

// type User struct {
// 	gorm.Model
// 	Role        string `json:"role"`
// 	FullName    string `json:"full_name"`
// 	Email       string `json:"email" gorm:"unique"`
// 	Password    string `json:"password"`
// 	Skills      string `json:"skills" gorm:"type:jsonb"`
// 	Country     string `json:"country"`
// 	City        string `json:"city"`
// 	DateOfBirth string `json:"date_of_birth"`
// 	PhoneNumber string `json:"phone_number" gorm:"unique"`
// }
