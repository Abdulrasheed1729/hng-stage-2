package service

import (
	"errors"
	"hng-stage2/internal/models"
	"hng-stage2/internal/repository"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository         repository.UserRepository
	organisationRepository repository.OrganisationRepository
}

func NewUserService(
	userRepository repository.UserRepository,
	organisationRepository repository.OrganisationRepository) *UserService {
	return &UserService{userRepository, organisationRepository}
}

func (s *UserService) Register(user *models.User) (*models.User, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	err = s.userRepository.Create(user)
	if err != nil {
		return nil, err
	}

	orgId := uuid.New()

	org := &models.Organisation{
		OrgID:       orgId.String(),
		Name:        user.FirstName + "'s Organisation",
		Description: "Default organisation for " + user.FirstName,
		// Users:       []models.User{*user},
	}

	err = s.organisationRepository.Create(org)
	log.Println(org)
	if err != nil {
		return nil, err
	}

	isUserInOrg, err := s.organisationRepository.IsUserInOrganisation(org.OrgID, user.UserID)

	if err != nil {

		// panic(err)
		return nil, err
	}

	if !isUserInOrg {
		// s.organisationRepository.AddUserToOrganisation(org.OrgID, user.UserID)

		org.Users = append(org.Users, *user)

		log.Println(org)

	}

	return user, nil
}

func (s *UserService) IsUserExisting(user models.User) (bool, error) {
	return s.userRepository.IsUserExisting(user)
}

func (s *UserService) Login(email, password string) (*models.User, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepository.GetUserByEmail(email)
}

func (s *UserService) ValidatePassword(user models.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
