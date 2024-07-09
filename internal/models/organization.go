package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organisation struct {
	gorm.Model
	OrgID       string `gorm:"primary_key;autoIncrement"` // orgId
	Name        string `gorm:"not null"`                  // name
	Description string `json:"description,omitempty"`     // description
}

func (org *Organisation) BeforeCreate(tx *gorm.DB) (err error) {
	org.OrgID = uuid.New().String()
	return
}

// func CreateOrg(organisation *Organisation) error {
// 	return database.Database.Create(organisation).Error
// }
// func GetUserOrgsByID(id string) ([]Organisation, error) {

// 	var orgs []Organisation

// 	err := database.Database.
// 		Preload("Users").
// 		Joins("JOIN user_organisations ON user_organisations.org_id = organisations.org_id").
// 		Joins("JOIN users ON userOrganisations.user_id = users.user_id").
// 		Where("users.user_id = ?", id).
// 		Find(&orgs).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return orgs, nil
// }

// func AddUserToOrganisation(orgID, userID string) error {
// 	// Add user to organisation
// 	err := database.Database.Exec("INSERT INTO user_organisations (org_id, user_id) VALUES (?, ?)", orgID, userID).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // func IsUserInOrganisation(orgID, userID string) (bool, error) {
// // 	var org Organisation
// // 	err := database.Database.
// // 		Where("org_id = ?", orgID).
// // 		Preload("Users").
// // 		First(&org).Error
// // 	if err != nil {
// // 		return false, err
// // 	}

// // 	for _, user := range org.Users {
// // 		if user.ID.String() == userID {
// // 			return true, nil
// // 		}
// // 	}
// // 	return false, nil
// // }

// func GetByOrgID(id string) (*Organisation, error) {
// 	var org Organisation

// 	err := database.Database.Where("org_id = ?", id).First(&org).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &org, nil
// }
