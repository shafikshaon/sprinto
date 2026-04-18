package repository

import (
	"gorm.io/gorm"
	"sprinto/models"
)

type DeadlineRepository interface {
	All(projectID uint) ([]models.Deadline, error)
	Create(d models.Deadline) error
	Update(id uint, title, project, dueDateRaw, priority string) error
	Delete(id uint) error
}

type deadlineRepo struct{ db *gorm.DB }

func NewDeadlineRepository(db *gorm.DB) DeadlineRepository { return &deadlineRepo{db: db} }

func (r *deadlineRepo) All(projectID uint) ([]models.Deadline, error) {
	var dl []models.Deadline
	q := r.db.Order("due_date ASC")
	if projectID > 0 {
		q = q.Where("project_id = ?", projectID)
	}
	return dl, q.Find(&dl).Error
}

func (r *deadlineRepo) Create(d models.Deadline) error { return r.db.Create(&d).Error }

func (r *deadlineRepo) Update(id uint, title, project, dueDateRaw, priority string) error {
	return r.db.Model(&models.Deadline{}).Where("id = ?", id).Updates(map[string]interface{}{
		"title":    title,
		"project":  project,
		"due_date": dueDateRaw,
		"priority": priority,
	}).Error
}

func (r *deadlineRepo) Delete(id uint) error { return r.db.Delete(&models.Deadline{}, id).Error }
