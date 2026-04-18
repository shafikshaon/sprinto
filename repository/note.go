package repository

import (
	"gorm.io/gorm"
	"sprinto/models"
)

type StickyNoteRepository interface {
	All(filter string) ([]models.StickyNote, error)
	GetByID(id uint) (models.StickyNote, error)
	Create(n models.StickyNote) error
	Update(id uint, title, content, color string) error
	TogglePin(id uint, pinned bool) error
	Delete(id uint) error
}

type stickyNoteRepo struct{ db *gorm.DB }

func NewStickyNoteRepository(db *gorm.DB) StickyNoteRepository {
	return &stickyNoteRepo{db: db}
}

func (r *stickyNoteRepo) All(filter string) ([]models.StickyNote, error) {
	var notes []models.StickyNote
	q := r.db.Order("pinned DESC, updated_at DESC")
	if filter == "pinned" {
		q = q.Where("pinned = ?", true)
	}
	return notes, q.Find(&notes).Error
}

func (r *stickyNoteRepo) GetByID(id uint) (models.StickyNote, error) {
	var n models.StickyNote
	return n, r.db.First(&n, id).Error
}

func (r *stickyNoteRepo) Create(n models.StickyNote) error { return r.db.Create(&n).Error }

func (r *stickyNoteRepo) Update(id uint, title, content, color string) error {
	return r.db.Model(&models.StickyNote{}).Where("id = ?", id).Updates(map[string]interface{}{
		"title":   title,
		"content": content,
		"color":   color,
	}).Error
}

func (r *stickyNoteRepo) TogglePin(id uint, pinned bool) error {
	return r.db.Model(&models.StickyNote{}).Where("id = ?", id).Update("pinned", pinned).Error
}

func (r *stickyNoteRepo) Delete(id uint) error {
	return r.db.Delete(&models.StickyNote{}, id).Error
}
