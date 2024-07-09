package repository

import (
	"hng-stage2/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
}
