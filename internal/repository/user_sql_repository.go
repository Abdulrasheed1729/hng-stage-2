package repository

import (
	"hng-stage2/internal/models"

	"gorm.io/gorm"
)

type UserSQLRepository struct {
	DB *gorm.DB
}

func NewUserSQLRepository(db *gorm.DB) UserRepository {
	return &UserSQLRepository{db}
}

func (r *UserSQLRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r UserSQLRepository) IsUserExisting(user models.User) (bool, error) {
	result := r.DB.Where("email = ? ", user.Email).Limit(1).Find(&user)

	if err := result.Error; err != nil {
		return false, err
	}

	if result.RowsAffected > 0 {
		return true, nil
	}

	return false, nil
}

func (r *UserSQLRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserSQLRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("user_id = ?", id).First(&user).Error
	return &user, err
}
