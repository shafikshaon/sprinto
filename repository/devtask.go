package repository

import (
	"gorm.io/gorm"
	"sprinto/models"
)

type DevTaskRepository interface {
	All(projectID uint) ([]models.DevTask, error)
	ByID(id uint) (models.DevTask, error)
	Create(t models.DevTask) error
	Update(id uint, title, typ, assignees, status, priority string) error
	Delete(id uint) error
	OpenCountsByType(projectID uint) (map[string]int, error)
	AddComment(c models.DevTaskComment) error
	DeleteComment(id uint) error
}

type devTaskRepo struct{ db *gorm.DB }

func NewDevTaskRepository(db *gorm.DB) DevTaskRepository { return &devTaskRepo{db: db} }

func (r *devTaskRepo) All(projectID uint) ([]models.DevTask, error) {
	var tasks []models.DevTask
	q := r.db.Preload("Comments").Order("created_at DESC")
	if projectID > 0 {
		q = q.Where("project_id = ?", projectID)
	}
	return tasks, q.Find(&tasks).Error
}

func (r *devTaskRepo) ByID(id uint) (models.DevTask, error) {
	var t models.DevTask
	result := r.db.Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).First(&t, id)
	return t, result.Error
}

func (r *devTaskRepo) Create(t models.DevTask) error { return r.db.Create(&t).Error }

func (r *devTaskRepo) Update(id uint, title, typ, assignees, status, priority string) error {
	return r.db.Model(&models.DevTask{}).Where("id = ?", id).
		Updates(map[string]interface{}{"title": title, "type": typ, "assignees": assignees, "status": status, "priority": priority}).Error
}

func (r *devTaskRepo) Delete(id uint) error { return r.db.Delete(&models.DevTask{}, id).Error }

func (r *devTaskRepo) AddComment(c models.DevTaskComment) error { return r.db.Create(&c).Error }

func (r *devTaskRepo) DeleteComment(id uint) error {
	return r.db.Delete(&models.DevTaskComment{}, id).Error
}

func (r *devTaskRepo) OpenCountsByType(projectID uint) (map[string]int, error) {
	type row struct {
		Type  string
		Count int
	}
	var rows []row
	q := r.db.Model(&models.DevTask{}).
		Select("type, COUNT(*) as count").
		Where("status != ?", "Done")
	if projectID > 0 {
		q = q.Where("project_id = ?", projectID)
	}
	if err := q.Group("type").Scan(&rows).Error; err != nil {
		return nil, err
	}
	counts := map[string]int{"Improvement": 0, "Tech Debt": 0, "Research": 0}
	for _, r := range rows {
		counts[r.Type] = r.Count
	}
	return counts, nil
}
