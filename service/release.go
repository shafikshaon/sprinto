package service

import (
	"strings"
	"time"

	"sprinto/models"
	"sprinto/repository"
)

type ReleaseService interface {
	All(projectID uint) ([]models.Release, error)
	ByID(id uint) (models.Release, error)
	Create(name, description, status, targetDate string, projectID uint) error
	Delete(id uint) error
	AddStage(releaseID uint, name, status string) error
	DeleteStage(id uint) error
	UpdateStageStatus(id uint, status string) error
	AddStory(stageID uint, title string, assigneeID *uint) error
	DeleteStory(id uint) error
	UpdateStoryStatus(id uint, status string) error
	UpdateStory(id uint, title string, assigneeID *uint) error
	AddSlackUpdate(stageID uint, channel, message, author string) error
	DeleteSlackUpdate(id uint) error
}

type releaseService struct{ repo repository.ReleaseRepository }

func NewReleaseService(r repository.ReleaseRepository) ReleaseService {
	return &releaseService{repo: r}
}

func (s *releaseService) All(projectID uint) ([]models.Release, error) {
	return s.repo.All(projectID)
}

func (s *releaseService) ByID(id uint) (models.Release, error) { return s.repo.ByID(id) }

func (s *releaseService) Create(name, description, status, targetDate string, projectID uint) error {
	if strings.TrimSpace(name) == "" {
		return nil
	}
	return s.repo.Create(models.Release{
		ProjectID:   projectID,
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		Status:      status,
		TargetDate:  targetDate,
	})
}

func (s *releaseService) Delete(id uint) error { return s.repo.Delete(id) }

func (s *releaseService) AddStage(releaseID uint, name, status string) error {
	if strings.TrimSpace(name) == "" {
		return nil
	}
	return s.repo.CreateStage(models.ReleaseStage{
		ReleaseID: releaseID,
		Name:      strings.TrimSpace(name),
		Status:    status,
	})
}

func (s *releaseService) DeleteStage(id uint) error        { return s.repo.DeleteStage(id) }
func (s *releaseService) UpdateStageStatus(id uint, status string) error {
	return s.repo.UpdateStageStatus(id, status)
}

func (s *releaseService) AddStory(stageID uint, title string, assigneeID *uint) error {
	if strings.TrimSpace(title) == "" {
		return nil
	}
	return s.repo.CreateStory(models.ReleaseStory{
		StageID:    stageID,
		Title:      strings.TrimSpace(title),
		AssigneeID: assigneeID,
		Status:     "Pending",
	})
}

func (s *releaseService) DeleteStory(id uint) error { return s.repo.DeleteStory(id) }
func (s *releaseService) UpdateStoryStatus(id uint, status string) error {
	return s.repo.UpdateStoryStatus(id, status)
}

func (s *releaseService) UpdateStory(id uint, title string, assigneeID *uint) error {
	if strings.TrimSpace(title) == "" {
		return nil
	}
	return s.repo.UpdateStory(id, strings.TrimSpace(title), assigneeID)
}

func (s *releaseService) AddSlackUpdate(stageID uint, channel, message, author string) error {
	if strings.TrimSpace(message) == "" {
		return nil
	}
	return s.repo.CreateSlackUpdate(models.ReleaseSlackUpdate{
		StageID:  stageID,
		Channel:  strings.TrimSpace(channel),
		Message:  strings.TrimSpace(message),
		Author:   strings.TrimSpace(author),
		PostedAt: time.Now().Format("Jan 2, 3:04 PM"),
	})
}

func (s *releaseService) DeleteSlackUpdate(id uint) error { return s.repo.DeleteSlackUpdate(id) }
