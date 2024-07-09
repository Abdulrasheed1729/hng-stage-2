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
	err := r.db.Where("org_id = ?", id).First(&org).Error
	return &org, err
}

func (r *OrganisationSQLRepository) GetUserOrganisations(userID string) ([]models.Organisation, error) {
	var orgs []models.Organisation
	// id, err := uuid.Parse(userID)

	// if err != nil {
	// 	return nil, err
	// }

	err := r.db.Model(&models.User{UserID: userID}).Association("Organisations").Find(&orgs)
	return orgs, err
}

func (r *OrganisationSQLRepository) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("user_id = ?", userID).First(&user).Error
	return &user, err
}

func (r *OrganisationSQLRepository) AddUserToOrganisation(orgID, userID string) error {
	return r.db.Exec("INSERT INTO Organisations (user_id, org_id) VALUES (?, ?)", userID, orgID).Error
}
