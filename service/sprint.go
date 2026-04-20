package service

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"sprinto/models"
	"sprinto/repository"
)

type SprintService interface {
	ActiveSprint(projectID uint) (models.Sprint, error)
	TaskByID(id uint) (models.SprintTask, error)
	AddTask(sprintID uint, title string, assigneeIDs []uint, status, priority string) error
	UpdateTask(id uint, title string, assigneeIDs []uint, status, priority string) error
	RemoveTask(id uint) error
	UpdateProgress(sprintID uint, progress int) error
	AddComment(taskID uint, author, content string) error
	DeleteComment(id uint) error
}

type sprintService struct{ repo repository.SprintRepository }

func NewSprintService(r repository.SprintRepository) SprintService {
	return &sprintService{repo: r}
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
	return s.repo.CreateTask(models.SprintTask{
		SprintID: sprintID,
		Title:    strings.TrimSpace(title),
		Status:   status,
		Priority: priority,
	}, assigneeIDs)
}

func (s *sprintService) TaskByID(id uint) (models.SprintTask, error) {
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
	return s.repo.AddComment(models.SprintTaskComment{
		TaskID:  taskID,
		Author:  author,
		Content: content,
	})
}

func (s *sprintService) DeleteComment(id uint) error { return s.repo.DeleteComment(id) }
