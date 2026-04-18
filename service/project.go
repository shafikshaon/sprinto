package service

import (
	"strings"

	"sprinto/models"
	"sprinto/repository"
)

type ProjectService interface {
	All() ([]models.Project, error)
	AllWithMembers() ([]models.Project, error)
	Create(name, description string) error
	Update(id uint, name, description string) error
	Delete(id uint) error
	AddMember(projectID, memberID uint) error
	RemoveMember(projectID, memberID uint) error
}

type projectService struct{ repo repository.ProjectRepository }

func NewProjectService(r repository.ProjectRepository) ProjectService {
	return &projectService{repo: r}
}

func (s *projectService) All() ([]models.Project, error)            { return s.repo.All() }
func (s *projectService) AllWithMembers() ([]models.Project, error) { return s.repo.AllWithMembers() }

func (s *projectService) Create(name, description string) error {
	if strings.TrimSpace(name) == "" {
		return nil
	}
	return s.repo.Create(models.Project{
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
	})
}

func (s *projectService) Update(id uint, name, description string) error {
	return s.repo.Update(id, strings.TrimSpace(name), strings.TrimSpace(description))
}

func (s *projectService) Delete(id uint) error          { return s.repo.Delete(id) }
func (s *projectService) AddMember(pid, mid uint) error  { return s.repo.AddMember(pid, mid) }
func (s *projectService) RemoveMember(pid, mid uint) error { return s.repo.RemoveMember(pid, mid) }
