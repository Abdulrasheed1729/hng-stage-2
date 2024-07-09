package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID    string `gorm:"primary_key" json:"userId"`
	FirstName string `gorm:"not null" json:"firstName"`
	LastName  string `gorm:"not null" json:"lastName"`
	Email     string `gorm:"unique;not null" json:"email"`
	Password  string `gorm:"not null" json:"password"`
	Phone     string `json:"phone"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.UserID = uuid.New().String()
	return
}

// func (user *User) Save() (*User, error) {
// 	err := database.Database.Create(&user).Error
// 	if err != nil {
// 		return &User{}, err
// 	}
// 	return user, nil
// }

// func (user *User) BeforeSave(*gorm.DB) error {

// 	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return err
// 	}
// 	user.Password = string(passwordHash)
// 	user.FirstName = html.EscapeString(strings.TrimSpace(user.FirstName))
// 	return nil
// }

// func FindUserByEmail(email string) (*User, error) {
// 	var user User

// 	err := database.Database.Where("email=?", email).Find(&user).Error

// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user, nil

// }

// func GetUserByID(id string) (*User, error) {
// 	var user User

// 	log.Printf("Querying user with ID: %s", id)

// 	err := database.Database.Where("user_id = ?", id).First(&user).Error
// 	if err != nil {
// 		log.Printf("Error fetching user by ID: %v", err)
// 		return nil, err
// 	}
// 	log.Printf("User found: %+v", user)

// 	return &user, nil
// }

// func (user *User) Create() (*User, error) {
// 	err := database.Database.Create(&user).Error
// 	if err != nil {
// 		return &User{}, err
// 	}
// 	return user, nil
// }