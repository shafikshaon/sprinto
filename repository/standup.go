package repository

import (
	"gorm.io/gorm"
	"sprinto/models"
)

type StandupFilter struct {
	DateFrom string
	DateTo   string
	Search   string
}

type StandupRepository interface {
	All(projectID uint, f StandupFilter, page, perPage int) ([]models.StandupEntry, int64, error)
	ByDate(date string, projectID uint) ([]models.StandupEntry, error)
	Create(e models.StandupEntry) error
	Update(id uint, summary, dependencies, issues, actionItems string, projectID uint) error
	Delete(id uint) error
}

type standupRepo struct{ db *gorm.DB }

func NewStandupRepository(db *gorm.DB) StandupRepository { return &standupRepo{db: db} }

func (r *standupRepo) filtered(projectID uint, f StandupFilter) *gorm.DB {
	q := r.db.Model(&models.StandupEntry{})
	if projectID > 0 {
		q = q.Where("project_id = ?", projectID)
	}
	if f.DateFrom != "" {
		q = q.Where("date >= ?", f.DateFrom)
	}
	if f.DateTo != "" {
		q = q.Where("date <= ?", f.DateTo)
	}
	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where("summary LIKE ? OR dependencies LIKE ? OR issues LIKE ? OR action_items LIKE ?", like, like, like, like)
	}
	return q
}

func (r *standupRepo) All(projectID uint, f StandupFilter, page, perPage int) ([]models.StandupEntry, int64, error) {
	var entries []models.StandupEntry
	var total int64
	if err := r.filtered(projectID, f).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	err := r.filtered(projectID, f).
		Preload("Project").
		Order("date DESC, created_at DESC").
		Offset(offset).Limit(perPage).
		Find(&entries).Error
	return entries, total, err
}

func (r *standupRepo) ByDate(date string, projectID uint) ([]models.StandupEntry, error) {
	var entries []models.StandupEntry
	q := r.db.Preload("Project").Where("date = ?", date)
	if projectID > 0 {
		q = q.Where("project_id = ?", projectID)
	}
	return entries, q.Order("created_at ASC").Find(&entries).Error
}

func (r *standupRepo) Create(e models.StandupEntry) error { return r.db.Create(&e).Error }

func (r *standupRepo) Update(id uint, summary, dependencies, issues, actionItems string, projectID uint) error {
	return r.db.Model(&models.StandupEntry{}).Where("id = ?", id).Updates(map[string]interface{}{
		"project_id":   projectID,
		"summary":      summary,
		"dependencies": dependencies,
		"issues":       issues,
		"action_items": actionItems,
	}).Error
}

func (r *standupRepo) Delete(id uint) error {
	return r.db.Delete(&models.StandupEntry{}, id).Error
}
