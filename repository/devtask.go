package repository

import (
	"gorm.io/gorm"
	"sprinto/models"
)

type DevTaskFilter struct {
	Search   string
	Type     string
	Status   string
	Priority string
}

type DevTaskRepository interface {
	All(projectID uint, f DevTaskFilter, page, perPage int) ([]models.DevTask, int64, error)
	ByID(id uint) (models.DevTask, error)
	Create(t models.DevTask, assigneeIDs []uint) error
	Update(id uint, title, typ, status, priority string, assigneeIDs []uint) error
	Delete(id uint) error
	OpenCountsByType(projectID uint) (map[string]int, error)
	AddComment(c models.DevTaskComment) error
	DeleteComment(id uint) error
}

type devTaskRepo struct{ db *gorm.DB }

func NewDevTaskRepository(db *gorm.DB) DevTaskRepository { return &devTaskRepo{db: db} }

func (r *devTaskRepo) filtered(projectID uint, f DevTaskFilter) *gorm.DB {
	q := r.db.Model(&models.DevTask{})
	if projectID > 0 {
		q = q.Where("project_id = ?", projectID)
	}
	if f.Search != "" {
		q = q.Where("title LIKE ?", "%"+f.Search+"%")
	}
	if f.Type != "" {
		q = q.Where("type = ?", f.Type)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.Priority != "" {
		q = q.Where("priority = ?", f.Priority)
	}
	return q
}

func (r *devTaskRepo) All(projectID uint, f DevTaskFilter, page, perPage int) ([]models.DevTask, int64, error) {
	var tasks []models.DevTask
	var total int64
	if err := r.filtered(projectID, f).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	err := r.filtered(projectID, f).
		Preload("Comments").
		Preload("Assignees").
		Order("created_at DESC").
		Offset(offset).Limit(perPage).
		Find(&tasks).Error
	return tasks, total, err
}

func (r *devTaskRepo) ByID(id uint) (models.DevTask, error) {
	var t models.DevTask
	result := r.db.Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).Preload("Assignees").First(&t, id)
	return t, result.Error
}

func (r *devTaskRepo) Create(t models.DevTask, assigneeIDs []uint) error {
	t.Assignees = nil
	if err := r.db.Create(&t).Error; err != nil {
		return err
	}
	if len(assigneeIDs) > 0 {
		members := make([]models.TeamMember, len(assigneeIDs))
		for i, id := range assigneeIDs {
			members[i] = models.TeamMember{Model: gorm.Model{ID: id}}
		}
		return r.db.Model(&t).Association("Assignees").Replace(members)
	}
	return nil
}

func (r *devTaskRepo) Update(id uint, title, typ, status, priority string, assigneeIDs []uint) error {
	if err := r.db.Model(&models.DevTask{}).Where("id = ?", id).
		Updates(map[string]interface{}{"title": title, "type": typ, "status": status, "priority": priority}).Error; err != nil {
		return err
	}
	task := models.DevTask{Model: gorm.Model{ID: id}}
	members := make([]models.TeamMember, len(assigneeIDs))
	for i, aid := range assigneeIDs {
		members[i] = models.TeamMember{Model: gorm.Model{ID: aid}}
	}
	return r.db.Model(&task).Association("Assignees").Replace(members)
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
