package repository

import (
	"gorm.io/gorm"
	"sprinto/models"
)

type ProjectRepository interface {
	All() ([]models.Project, error)
	AllWithMembers() ([]models.Project, error)
	ByID(id uint) (models.Project, error)
	Create(p models.Project) error
	Update(id uint, name, description string) error
	Delete(id uint) error
	AddMember(projectID, memberID uint) error
	RemoveMember(projectID, memberID uint) error
}

type projectRepo struct{ db *gorm.DB }

func NewProjectRepository(db *gorm.DB) ProjectRepository { return &projectRepo{db: db} }

func (r *projectRepo) All() ([]models.Project, error) {
	var projects []models.Project
	return projects, r.db.Order("created_at ASC").Find(&projects).Error
}

func (r *projectRepo) AllWithMembers() ([]models.Project, error) {
	var projects []models.Project
	return projects, r.db.Preload("Members").Order("created_at ASC").Find(&projects).Error
}

func (r *projectRepo) ByID(id uint) (models.Project, error) {
	var p models.Project
	return p, r.db.Preload("Members").First(&p, id).Error
}

func (r *projectRepo) Create(p models.Project) error { return r.db.Create(&p).Error }

func (r *projectRepo) Update(id uint, name, description string) error {
	return r.db.Model(&models.Project{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":        name,
		"description": description,
	}).Error
}

func (r *projectRepo) Delete(id uint) error {
	var p models.Project
	if err := r.db.First(&p, id).Error; err != nil {
		return err
	}
	r.db.Model(&p).Association("Members").Clear()
	return r.db.Delete(&p).Error
}

func (r *projectRepo) AddMember(projectID, memberID uint) error {
	var p models.Project
	if err := r.db.First(&p, projectID).Error; err != nil {
		return err
	}
	var m models.TeamMember
	if err := r.db.First(&m, memberID).Error; err != nil {
		return err
	}
	return r.db.Model(&p).Association("Members").Append(&m)
}

func (r *projectRepo) RemoveMember(projectID, memberID uint) error {
	var p models.Project
	if err := r.db.First(&p, projectID).Error; err != nil {
		return err
	}
	var m models.TeamMember
	if err := r.db.First(&m, memberID).Error; err != nil {
		return err
	}
	return r.db.Model(&p).Association("Members").Delete(&m)
}
