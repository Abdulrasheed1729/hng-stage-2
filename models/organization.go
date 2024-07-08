package models

import "gorm.io/gorm"

type Organisation struct {
	gorm.Model
	ID          string `gorm:"primary_key;unique"`    // orgId
	Name        string `gorm:"not null"`              // name
	Description string `json:"description,omitempty"` // description
	Users       []User
}
