package repository

import "hng-stage2/internal/models"

type OrganisationRepository interface {
	Create(org *models.Organisation) error
	GetOrganisationByID(id string) (*models.Organisation, error)
	GetUserOrganisations(userID string) ([]models.Organisation, error)
	GetUserByID(userID string) (*models.User, error)
	AddUserToOrganisation(orgID, userID string) error
	IsUserInOrganisation(orgID, userID string) (bool, error)
}
