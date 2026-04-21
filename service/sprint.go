package service

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"sprinto/models"
	"sprinto/repository"
)

type SprintService interface {
	// Sprint management
	ListSprints(projectID uint) ([]models.Sprint, error)
	CreateSprint(projectID uint, name, goal, startDate, endDate string) error
	UpdateSprint(id uint, name, goal, startDate, endDate string) error
	DeleteSprint(id uint) error
	ActivateSprint(id uint, projectID uint) error
	// Board
	ActiveSprint(projectID uint) (models.Sprint, error)
	TaskByID(id uint) (models.Task, error)
	AddTask(sprintID uint, title string, assigneeIDs []uint, status, priority string) error
	UpdateTask(id uint, title string, assigneeIDs []uint, status, priority string) error
	RemoveTask(id uint) error
	UpdateProgress(sprintID uint, progress int) error
	AddComment(taskID uint, author, content string) error
	DeleteComment(id uint) error
	// Release (merged)
	UpdateRelease(id uint, description, status, targetDate string) error
	AddStage(sprintID uint, name, status string) error
	DeleteStage(id uint) error
	UpdateStageStatus(id uint, status string) error
	AddStory(stageID uint, title string, assigneeID *uint) error
	DeleteStory(id uint) error
	UpdateStoryStatus(id uint, status string) error
	UpdateStory(id uint, title string, assigneeID *uint) error
	AddSlackUpdate(stageID uint, channel, message, author string) error
	DeleteSlackUpdate(id uint) error
}

type sprintService struct{ repo repository.SprintRepository }

func NewSprintService(r repository.SprintRepository) SprintService {
	return &sprintService{repo: r}
}

func (s *sprintService) ListSprints(projectID uint) ([]models.Sprint, error) {
	return s.repo.ListSprints(projectID)
}

func (s *sprintService) CreateSprint(projectID uint, name, goal, startDate, endDate string) error {
	if strings.TrimSpace(name) == "" {
		return nil
	}
	return s.repo.CreateSprint(models.Sprint{
		ProjectID: projectID,
		Name:      strings.TrimSpace(name),
		Goal:      strings.TrimSpace(goal),
		StartDate: startDate,
		EndDate:   endDate,
		Status:    "Draft",
	})
}

func (s *sprintService) UpdateSprint(id uint, name, goal, startDate, endDate string) error {
	if strings.TrimSpace(name) == "" {
		return nil
	}
	return s.repo.UpdateSprint(id, strings.TrimSpace(name), strings.TrimSpace(goal), startDate, endDate)
}

func (s *sprintService) DeleteSprint(id uint) error { return s.repo.DeleteSprint(id) }

func (s *sprintService) ActivateSprint(id uint, projectID uint) error {
	return s.repo.ActivateSprint(id, projectID)
}

func (s *sprintService) ActiveSprint(projectID uint) (models.Sprint, error) {
	sprint, err := s.repo.ActiveSprint(projectID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Sprint{}, nil
	}
	return sprint, err
}

func (s *sprintService) AddTask(sprintID uint, title string, assigneeIDs []uint, status, priority string) error {
	if strings.TrimSpace(title) == "" {
		return nil
	}
	return s.repo.CreateTask(models.Task{
		Category: "sprint",
		SprintID: sprintID,
		Title:    strings.TrimSpace(title),
		Status:   status,
		Priority: priority,
	}, assigneeIDs)
}

func (s *sprintService) TaskByID(id uint) (models.Task, error) {
	return s.repo.TaskByID(id)
}

func (s *sprintService) UpdateTask(id uint, title string, assigneeIDs []uint, status, priority string) error {
	if strings.TrimSpace(title) == "" {
		return nil
	}
	return s.repo.UpdateTask(id, strings.TrimSpace(title), status, priority, assigneeIDs)
}

func (s *sprintService) RemoveTask(id uint) error { return s.repo.DeleteTask(id) }

func (s *sprintService) UpdateProgress(id uint, p int) error {
	return s.repo.UpdateProgress(id, clamp(p, 0, 100))
}

func (s *sprintService) AddComment(taskID uint, author, content string) error {
	author = strings.TrimSpace(author)
	content = strings.TrimSpace(content)
	if content == "" {
		return nil
	}
	if author == "" {
		author = "Anonymous"
	}
	return s.repo.AddComment(models.TaskComment{
		TaskID:  taskID,
		Author:  author,
		Content: content,
	})
}

func (s *sprintService) DeleteComment(id uint) error { return s.repo.DeleteComment(id) }

// ── Release methods ───────────────────────────────────────────────────────────

func (s *sprintService) UpdateRelease(id uint, description, status, targetDate string) error {
	return s.repo.UpdateRelease(id, strings.TrimSpace(description), status, targetDate)
}

func (s *sprintService) AddStage(sprintID uint, name, status string) error {
	if strings.TrimSpace(name) == "" {
		return nil
	}
	return s.repo.CreateStage(models.ReleaseStage{
		SprintID: sprintID,
		Name:     strings.TrimSpace(name),
		Status:   status,
	})
}

func (s *sprintService) DeleteStage(id uint) error { return s.repo.DeleteStage(id) }

func (s *sprintService) UpdateStageStatus(id uint, status string) error {
	return s.repo.UpdateStageStatus(id, status)
}

func (s *sprintService) AddStory(stageID uint, title string, assigneeID *uint) error {
	if strings.TrimSpace(title) == "" {
		return nil
	}
	return s.repo.CreateStory(models.Task{
		Category:       "release",
		ReleaseStageID: stageID,
		Title:          strings.TrimSpace(title),
		AssigneeID:     assigneeID,
		Status:         "Pending",
	})
}

func (s *sprintService) DeleteStory(id uint) error { return s.repo.DeleteStory(id) }

func (s *sprintService) UpdateStoryStatus(id uint, status string) error {
	return s.repo.UpdateStoryStatus(id, status)
}

func (s *sprintService) UpdateStory(id uint, title string, assigneeID *uint) error {
	if strings.TrimSpace(title) == "" {
		return nil
	}
	return s.repo.UpdateStory(id, strings.TrimSpace(title), assigneeID)
}

func (s *sprintService) AddSlackUpdate(stageID uint, channel, message, author string) error {
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

func (s *sprintService) DeleteSlackUpdate(id uint) error { return s.repo.DeleteSlackUpdate(id) }
