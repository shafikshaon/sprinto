package service

import (
	"strings"

	"sprinto/models"
	"sprinto/repository"
)

type StickyNoteService interface {
	All(filter string) ([]models.StickyNote, error)
	GetByID(id uint) (models.StickyNote, error)
	Create(title, content, color string) error
	Update(id uint, title, content, color string) error
	TogglePin(id uint, pinned bool) error
	Delete(id uint) error
}

type stickyNoteService struct{ repo repository.StickyNoteRepository }

func NewStickyNoteService(r repository.StickyNoteRepository) StickyNoteService {
	return &stickyNoteService{repo: r}
}

func (s *stickyNoteService) All(filter string) ([]models.StickyNote, error) {
	return s.repo.All(filter)
}

func (s *stickyNoteService) GetByID(id uint) (models.StickyNote, error) {
	return s.repo.GetByID(id)
}

func (s *stickyNoteService) Create(title, content, color string) error {
	if strings.TrimSpace(content) == "" {
		return nil
	}
	if color == "" {
		color = "yellow"
	}
	return s.repo.Create(models.StickyNote{
		Title:   strings.TrimSpace(title),
		Content: content,
		Color:   color,
	})
}

func (s *stickyNoteService) Update(id uint, title, content, color string) error {
	if color == "" {
		color = "yellow"
	}
	return s.repo.Update(id, strings.TrimSpace(title), content, color)
}

func (s *stickyNoteService) TogglePin(id uint, pinned bool) error {
	return s.repo.TogglePin(id, pinned)
}

func (s *stickyNoteService) Delete(id uint) error { return s.repo.Delete(id) }
