package service

import (
	"strings"

	"sprinto/models"
	"sprinto/repository"
)

type DevTaskService interface {
	All(projectID uint, f repository.DevTaskFilter, page, perPage int) ([]models.DevTask, int64, error)
	ByID(id uint) (models.DevTask, error)
	Add(title, typ string, assigneeIDs []uint, status, priority string, projectID uint) error
	Update(id uint, title, typ string, assigneeIDs []uint, status, priority string) error
	Remove(id uint) error
	OpenCountsByType(projectID uint) (map[string]int, error)
	AddComment(taskID uint, author, content string) error
	DeleteComment(id uint) error
}

type devTaskService struct{ repo repository.DevTaskRepository }

func NewDevTaskService(r repository.DevTaskRepository) DevTaskService {
	return &devTaskService{repo: r}
}

func (s *devTaskService) All(projectID uint, f repository.DevTaskFilter, page, perPage int) ([]models.DevTask, int64, error) {
	return s.repo.All(projectID, f, page, perPage)
}

func (s *devTaskService) ByID(id uint) (models.DevTask, error) {
	return s.repo.ByID(id)
}

func (s *devTaskService) Add(title, typ string, assigneeIDs []uint, status, priority string, projectID uint) error {
	if strings.TrimSpace(title) == "" {
		return nil
	}
	return s.repo.Create(models.DevTask{
		ProjectID: projectID,
		Title:     strings.TrimSpace(title),
		Type:      typ,
		Status:    status,
		Priority:  priority,
	}, assigneeIDs)
}

func (s *devTaskService) Update(id uint, title, typ string, assigneeIDs []uint, status, priority string) error {
	if strings.TrimSpace(title) == "" {
		return nil
	}
	return s.repo.Update(id, strings.TrimSpace(title), typ, status, priority, assigneeIDs)
}

func (s *devTaskService) Remove(id uint) error { return s.repo.Delete(id) }

func (s *devTaskService) OpenCountsByType(projectID uint) (map[string]int, error) {
	return s.repo.OpenCountsByType(projectID)
}

func (s *devTaskService) AddComment(taskID uint, author, content string) error {
	author = strings.TrimSpace(author)
	content = strings.TrimSpace(content)
	if content == "" {
		return nil
	}
	if author == "" {
		author = "Anonymous"
	}
	return s.repo.AddComment(models.DevTaskComment{
		TaskID:  taskID,
		Author:  author,
		Content: content,
	})
}

func (s *devTaskService) DeleteComment(id uint) error { return s.repo.DeleteComment(id) }
