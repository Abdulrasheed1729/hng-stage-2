package service

import (
	"hng-stage2/internal/models"
	"hng-stage2/internal/repository"
)

type OrganisationService struct {
	organisationRepository repository.OrganisationRepository
}

func NewOrganisationService(organisationRepository repository.OrganisationRepository) *OrganisationService {
	return &OrganisationService{organisationRepository}
}

func (s *OrganisationService) Create(org *models.Organisation) error {
	return s.organisationRepository.Create(org)
}

func (s *OrganisationService) GetByID(id string) (*models.Organisation, error) {
	return s.organisationRepository.GetOrganisationByID(id)
}

func (s *OrganisationService) GetByUser(userID string) ([]models.Organisation, error) {
	return s.organisationRepository.GetUserOrganisations(userID)
}

func (s *OrganisationService) GetUserByID(userID string) (*models.User, error) {
	return s.organisationRepository.GetUserByID(userID)
}

func (s *OrganisationService) AddUserToOrganisation(orgID, userID string) error {
	return s.organisationRepository.AddUserToOrganisation(orgID, userID)
}
