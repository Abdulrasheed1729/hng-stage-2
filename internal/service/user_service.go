package service

import (
	"errors"
	"hng-stage2/internal/models"
	"hng-stage2/internal/repository"

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
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	err = s.userRepository.Create(user)
	if err != nil {
		return nil, err
	}

	org := &models.Organisation{
		Name: user.FirstName + "'s Organisation",
	}
	err = s.organisationRepository.Create(org)
	if err != nil {
		return nil, err
	}

	return user, nil
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
