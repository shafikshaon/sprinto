package service

import (
	"strings"

	"sprinto/models"
	"sprinto/repository"
)

type DeadlineService interface {
	All(projectID uint) ([]models.Deadline, error)
	Add(title, project, dueDateRaw, priority string, projectID uint) error
	Update(id uint, title, project, dueDateRaw, priority string) error
	Remove(id uint) error
}

type deadlineService struct{ repo repository.DeadlineRepository }

func NewDeadlineService(r repository.DeadlineRepository) DeadlineService {
	return &deadlineService{repo: r}
}

func (s *deadlineService) All(projectID uint) ([]models.Deadline, error) {
	deadlines, err := s.repo.All(projectID)
	if err != nil {
		return nil, err
	}
	for i := range deadlines {
		deadlines[i].DaysLeft = daysLeft(deadlines[i].DueDateRaw)
		deadlines[i].DueDate = formatDate(deadlines[i].DueDateRaw)
	}
	return deadlines, nil
}

func (s *deadlineService) Add(title, project, dueDateRaw, priority string, projectID uint) error {
	if strings.TrimSpace(title) == "" || dueDateRaw == "" {
		return nil
	}
	return s.repo.Create(models.Deadline{
		ProjectID:  projectID,
		Title:      strings.TrimSpace(title),
		Project:    strings.TrimSpace(project),
		DueDateRaw: dueDateRaw,
		Priority:   priority,
	})
}

func (s *deadlineService) Update(id uint, title, project, dueDateRaw, priority string) error {
	return s.repo.Update(id, strings.TrimSpace(title), strings.TrimSpace(project), dueDateRaw, priority)
}

func (s *deadlineService) Remove(id uint) error { return s.repo.Delete(id) }
