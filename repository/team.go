package repository

import (
	"gorm.io/gorm"
	"sprinto/models"
)

type TeamMemberRepository interface {
	All() ([]models.TeamMember, error)
	Create(m models.TeamMember) error
	Delete(id uint) error
}

type teamMemberRepo struct{ db *gorm.DB }

func NewTeamMemberRepository(db *gorm.DB) TeamMemberRepository {
	return &teamMemberRepo{db: db}
}

func (r *teamMemberRepo) All() ([]models.TeamMember, error) {
	var members []models.TeamMember
	return members, r.db.Preload("Projects").Order("created_at ASC").Find(&members).Error
}

func (r *teamMemberRepo) Create(m models.TeamMember) error { return r.db.Create(&m).Error }

func (r *teamMemberRepo) Delete(id uint) error {
	var m models.TeamMember
	if err := r.db.First(&m, id).Error; err != nil {
		return err
	}
	r.db.Model(&m).Association("Projects").Clear()
	return r.db.Delete(&m).Error
}
