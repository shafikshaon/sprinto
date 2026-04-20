package repository

import (
	"gorm.io/gorm"
	"sprinto/models"
)

type SprintRepository interface {
	ActiveSprint(projectID uint) (models.Sprint, error)
	TaskByID(id uint) (models.SprintTask, error)
	CreateTask(t models.SprintTask, assigneeIDs []uint) error
	UpdateTask(id uint, title, status, priority string, assigneeIDs []uint) error
	DeleteTask(id uint) error
	UpdateProgress(sprintID uint, progress int) error
	AddComment(c models.SprintTaskComment) error
	DeleteComment(id uint) error
}

type sprintRepo struct{ db *gorm.DB }

func NewSprintRepository(db *gorm.DB) SprintRepository { return &sprintRepo{db: db} }

func (r *sprintRepo) ActiveSprint(projectID uint) (models.Sprint, error) {
	var s models.Sprint
	q := r.db.Preload("Tasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).Preload("Tasks.Assignees").Preload("Tasks.Comments").Where("active = ?", true)
	if projectID > 0 {
		q = q.Where("project_id = ?", projectID)
	}
	result := q.Order("id DESC").First(&s)
	return s, result.Error
}

func (r *sprintRepo) TaskByID(id uint) (models.SprintTask, error) {
	var t models.SprintTask
	result := r.db.Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).Preload("Assignees").First(&t, id)
	return t, result.Error
}

func (r *sprintRepo) CreateTask(t models.SprintTask, assigneeIDs []uint) error {
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

func (r *sprintRepo) UpdateTask(id uint, title, status, priority string, assigneeIDs []uint) error {
	if err := r.db.Model(&models.SprintTask{}).Where("id = ?", id).
		Updates(map[string]interface{}{"title": title, "status": status, "priority": priority}).Error; err != nil {
		return err
	}
	task := models.SprintTask{Model: gorm.Model{ID: id}}
	members := make([]models.TeamMember, len(assigneeIDs))
	for i, aid := range assigneeIDs {
		members[i] = models.TeamMember{Model: gorm.Model{ID: aid}}
	}
	return r.db.Model(&task).Association("Assignees").Replace(members)
}

func (r *sprintRepo) DeleteTask(id uint) error { return r.db.Delete(&models.SprintTask{}, id).Error }

func (r *sprintRepo) UpdateProgress(sprintID uint, progress int) error {
	return r.db.Model(&models.Sprint{}).Where("id = ?", sprintID).Update("progress", progress).Error
}

func (r *sprintRepo) AddComment(c models.SprintTaskComment) error { return r.db.Create(&c).Error }

func (r *sprintRepo) DeleteComment(id uint) error {
	return r.db.Delete(&models.SprintTaskComment{}, id).Error
}
