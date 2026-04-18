package repository

import (
	"gorm.io/gorm"
	"sprinto/models"
)

type SprintRepository interface {
	ActiveSprint(projectID uint) (models.Sprint, error)
	TaskByID(id uint) (models.SprintTask, error)
	CreateTask(t models.SprintTask) error
	UpdateTask(id uint, title, assignees, status, priority string) error
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
	}).Preload("Tasks.Comments").Where("active = ?", true)
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
	}).First(&t, id)
	return t, result.Error
}

func (r *sprintRepo) CreateTask(t models.SprintTask) error { return r.db.Create(&t).Error }

func (r *sprintRepo) UpdateTask(id uint, title, assignees, status, priority string) error {
	return r.db.Model(&models.SprintTask{}).Where("id = ?", id).
		Updates(map[string]interface{}{"title": title, "assignees": assignees, "status": status, "priority": priority}).Error
}

func (r *sprintRepo) DeleteTask(id uint) error { return r.db.Delete(&models.SprintTask{}, id).Error }

func (r *sprintRepo) UpdateProgress(sprintID uint, progress int) error {
	return r.db.Model(&models.Sprint{}).Where("id = ?", sprintID).Update("progress", progress).Error
}

func (r *sprintRepo) AddComment(c models.SprintTaskComment) error { return r.db.Create(&c).Error }

func (r *sprintRepo) DeleteComment(id uint) error {
	return r.db.Delete(&models.SprintTaskComment{}, id).Error
}
