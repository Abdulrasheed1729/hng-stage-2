package models

import (
	"hng-stage2/database"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type User struct {
	gorm.Model
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"userId"`
	FirstName     string         `gorm:"not null" json:"firstName"`
	LastName      string         `gorm:"not null" json:"lastName"`
	Email         string         `gorm:"unique;not null" json:"email"`
	Password      string         `gorm:"not null" json:"password"`
	Phone         string         `gorm:"-" json:"phone"`
	Organizations []Organisation `gorm:"many2many:userOrganizations"`
}

func (user *User) BeforeSave(*gorm.DB) error {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	user.FirstName = html.EscapeString(strings.TrimSpace(user.FirstName))
	return nil
}

func FindUserByEmail(email string) (*User, error) {
	var user User

	err := database.Database.Where("email=?", email).Find(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil

}
