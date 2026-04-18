package service

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"sprinto/models"
	"sprinto/repository"
)

type AuthService interface {
	Register(fullName, email, password string) error
	Login(email, password string) (models.User, error)
	UserByID(id uint) (models.User, error)
}

type authService struct{ repo repository.UserRepository }

func NewAuthService(r repository.UserRepository) AuthService {
	return &authService{repo: r}
}

func (s *authService) Register(fullName, email, password string) error {
	fullName = strings.TrimSpace(fullName)
	email = strings.TrimSpace(email)
	if fullName == "" || email == "" || password == "" {
		return errors.New("all fields are required")
	}
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.Create(models.User{
		FullName:     fullName,
		Email:        email,
		PasswordHash: string(hash),
	})
}

func (s *authService) Login(email, password string) (models.User, error) {
	user, err := s.repo.ByEmail(strings.TrimSpace(email))
	if err != nil {
		return models.User{}, errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return models.User{}, errors.New("invalid email or password")
	}
	return user, nil
}

func (s *authService) UserByID(id uint) (models.User, error) {
	return s.repo.ByID(id)
}
