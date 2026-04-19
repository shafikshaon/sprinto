package service

import (
	"strings"

	"sprinto/models"
	"sprinto/repository"
)

type TeamMemberService interface {
	All() ([]models.TeamMember, error)
	Create(name, role, email string) error
	CreateForUser(name, role, email string, userID uint) error
	Delete(id uint) error
}

type teamMemberService struct{ repo repository.TeamMemberRepository }

func NewTeamMemberService(r repository.TeamMemberRepository) TeamMemberService {
	return &teamMemberService{repo: r}
}

func (s *teamMemberService) All() ([]models.TeamMember, error) { return s.repo.All() }

func (s *teamMemberService) Create(name, role, email string) error {
	if strings.TrimSpace(name) == "" {
		return nil
	}
	return s.repo.Create(models.TeamMember{
		Name:  strings.TrimSpace(name),
		Role:  strings.TrimSpace(role),
		Email: strings.TrimSpace(email),
	})
}

func (s *teamMemberService) CreateForUser(name, role, email string, userID uint) error {
	if strings.TrimSpace(name) == "" {
		return nil
	}
	return s.repo.Create(models.TeamMember{
		UserID: &userID,
		Name:   strings.TrimSpace(name),
		Role:   strings.TrimSpace(role),
		Email:  strings.TrimSpace(email),
	})
}

func (s *teamMemberService) Delete(id uint) error { return s.repo.Delete(id) }
