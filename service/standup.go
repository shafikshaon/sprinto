package service

import (
	"strings"
	"time"

	"sprinto/models"
	"sprinto/repository"
)

type StandupService interface {
	All(projectID uint, f repository.StandupFilter, page, perPage int) ([]models.StandupEntry, int64, error)
	ByDate(date string, projectID uint) ([]models.StandupEntry, error)
	Add(date, summary, dependencies, issues, actionItems string, projectID uint) error
	Update(id uint, summary, dependencies, issues, actionItems string, projectID uint) error
	Remove(id uint) error
}

type standupService struct{ repo repository.StandupRepository }

func NewStandupService(r repository.StandupRepository) StandupService {
	return &standupService{repo: r}
}

func (s *standupService) All(projectID uint, f repository.StandupFilter, page, perPage int) ([]models.StandupEntry, int64, error) {
	return s.repo.All(projectID, f, page, perPage)
}

func (s *standupService) ByDate(date string, projectID uint) ([]models.StandupEntry, error) {
	return s.repo.ByDate(date, projectID)
}

func (s *standupService) Add(date, summary, dependencies, issues, actionItems string, projectID uint) error {
	if strings.TrimSpace(summary) == "" && strings.TrimSpace(issues) == "" {
		return nil
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	return s.repo.Create(models.StandupEntry{
		ProjectID:    projectID,
		Date:         date,
		Summary:      strings.TrimSpace(summary),
		Dependencies: strings.TrimSpace(dependencies),
		Issues:       strings.TrimSpace(issues),
		ActionItems:  strings.TrimSpace(actionItems),
	})
}

func (s *standupService) Update(id uint, summary, dependencies, issues, actionItems string, projectID uint) error {
	return s.repo.Update(id,
		strings.TrimSpace(summary),
		strings.TrimSpace(dependencies),
		strings.TrimSpace(issues),
		strings.TrimSpace(actionItems),
		projectID,
	)
}

func (s *standupService) Remove(id uint) error { return s.repo.Delete(id) }
