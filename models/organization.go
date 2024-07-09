package models

import (
	"hng-stage2/database"

	"gorm.io/gorm"
)

type Organisation struct {
	gorm.Model
	ID          string `gorm:"primary_key;unique"`    // orgId
	Name        string `gorm:"not null"`              // name
	Description string `json:"description,omitempty"` // description
	Users       []User
}

func CreateOrg(organisation *Organisation) error {
	return database.Database.Create(organisation).Error
}
func GetUserOrgsByID(id string) ([]Organisation, error) {

	var orgs []Organisation

	err := database.Database.
		Preload("Users").
		Joins("JOIN userOrganisations ON userOrganisations.orgId = organisations.orgId").
		Joins("JOIN users ON userOrganisations.userId = users.userId").
		Where("users.userId = ?", id).
		Find(&orgs).Error
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func AddUserToOrganisation(orgID, userID string) error {
	// Add user to organisation
	err := database.Database.Exec("INSERT INTO userOrganisations (orgId, userId) VALUES (?, ?)", orgID, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func IsUserInOrganisation(orgID, userID string) (bool, error) {
	var org Organisation
	err := database.Database.
		Where("orgId = ?", orgID).
		Preload("Users").
		First(&org).Error
	if err != nil {
		return false, err
	}

	for _, user := range org.Users {
		if user.ID.String() == userID {
			return true, nil
		}
	}
	return false, nil
}

func GetByOrgID(id string) (*Organisation, error) {
	var org Organisation

	err := database.Database.Where("orgId = ?", id).First(&org).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}
