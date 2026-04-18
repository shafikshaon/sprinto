package repository

import (
	"gorm.io/gorm"
	"sprinto/models"
)

type StandupRepository interface {
	All(projectID uint) ([]models.StandupEntry, error)
	ByDate(date string, projectID uint) ([]models.StandupEntry, error)
	Create(e models.StandupEntry) error
	Update(id uint, summary, dependencies, issues, actionItems string) error
	Delete(id uint) error
}

type standupRepo struct{ db *gorm.DB }

func NewStandupRepository(db *gorm.DB) StandupRepository { return &standupRepo{db: db} }

func (r *standupRepo) All(projectID uint) ([]models.StandupEntry, error) {
	var entries []models.StandupEntry
	q := r.db.Order("date DESC, created_at DESC")
	if projectID > 0 {
		q = q.Where("project_id = ?", projectID)
	}
	return entries, q.Find(&entries).Error
}

func (r *standupRepo) ByDate(date string, projectID uint) ([]models.StandupEntry, error) {
	var entries []models.StandupEntry
	q := r.db.Where("date = ?", date)
	if projectID > 0 {
		q = q.Where("project_id = ?", projectID)
	}
	return entries, q.Order("created_at ASC").Find(&entries).Error
}

func (r *standupRepo) Create(e models.StandupEntry) error { return r.db.Create(&e).Error }

func (r *standupRepo) Update(id uint, summary, dependencies, issues, actionItems string) error {
	return r.db.Model(&models.StandupEntry{}).Where("id = ?", id).Updates(map[string]interface{}{
		"summary":      summary,
		"dependencies": dependencies,
		"issues":       issues,
		"action_items": actionItems,
	}).Error
}

func (r *standupRepo) Delete(id uint) error {
	return r.db.Delete(&models.StandupEntry{}, id).Error
}
