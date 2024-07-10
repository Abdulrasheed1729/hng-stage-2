package repository

import (
	"hng-stage2/internal/models"

	"gorm.io/gorm"
)

type OrganisationSQLRepository struct {
	db *gorm.DB
}

func NewOrganisationSQLRepository(db *gorm.DB) OrganisationRepository {
	return &OrganisationSQLRepository{db}
}

func (r *OrganisationSQLRepository) Create(org *models.Organisation) error {
	return r.db.Create(org).Error
}

func (r *OrganisationSQLRepository) GetOrganisationByID(id string) (*models.Organisation, error) {
	var org models.Organisation
	// err := m.DB.Preload("Users").Where("org_id = ?", id).First(&org).Error
	err := r.db.Where("org_id = ?", id).First(&org).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *OrganisationSQLRepository) GetUserOrganisations(userID string) ([]models.Organisation, error) {
	var orgs []models.Organisation

	// Adjust the query to use the correct column names from user_organisations table
	err := r.db.
		Preload("Users").
		Joins("JOIN user_organisations ON user_organisations.org_id = organisations.org_id").
		Joins("JOIN users ON user_organisations.user_id = users.user_id").
		Where("users.user_id = ?", userID).
		Find(&orgs).Error
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func (r *OrganisationSQLRepository) IsUserInOrganisation(orgID, userID string) (bool, error) {
	var org models.Organisation
	err := r.db.
		Where("org_id = ?", orgID).
		Preload("Users").
		First(&org).Error
	if err != nil {
		return false, err
	}

	for _, user := range org.Users {
		if user.UserID == userID {
			return true, nil
		}
	}
	return false, nil
}

func (r *OrganisationSQLRepository) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("user_id = ?", userID).First(&user).Error
	return &user, err
}

func (r *OrganisationSQLRepository) AddUserToOrganisation(orgID, userID string) error {

	err := r.db.Exec("INSERT INTO user_organisations (org_id, user_id) VALUES (?, ?)", orgID, userID).Error
	if err != nil {
		return err
	}
	return nil
}
