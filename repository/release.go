package repository

import (
	"gorm.io/gorm"
	"sprinto/models"
)

type ReleaseRepository interface {
	All(projectID uint) ([]models.Release, error)
	ByID(id uint) (models.Release, error)
	Create(r models.Release) error
	Delete(id uint) error
	CreateStage(s models.ReleaseStage) error
	DeleteStage(id uint) error
	UpdateStageStatus(id uint, status string) error
	CreateStory(s models.ReleaseStory) error
	DeleteStory(id uint) error
	UpdateStoryStatus(id uint, status string) error
	UpdateStory(id uint, title, assignee string) error
	CreateSlackUpdate(u models.ReleaseSlackUpdate) error
	DeleteSlackUpdate(id uint) error
}

type releaseRepo struct{ db *gorm.DB }

func NewReleaseRepository(db *gorm.DB) ReleaseRepository { return &releaseRepo{db: db} }

func (r *releaseRepo) All(projectID uint) ([]models.Release, error) {
	var releases []models.Release
	q := r.db.Preload("Stages").Order("created_at DESC")
	if projectID > 0 {
		q = q.Where("project_id = ?", projectID)
	}
	return releases, q.Find(&releases).Error
}

func (r *releaseRepo) ByID(id uint) (models.Release, error) {
	var rel models.Release
	result := r.db.
		Preload("Stages", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
		Preload("Stages.Stories", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
		Preload("Stages.SlackUpdates", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
		First(&rel, id)
	return rel, result.Error
}

func (r *releaseRepo) Create(rel models.Release) error { return r.db.Create(&rel).Error }

func (r *releaseRepo) Delete(id uint) error {
	var stageIDs []uint
	r.db.Model(&models.ReleaseStage{}).Where("release_id = ?", id).Pluck("id", &stageIDs)
	for _, sid := range stageIDs {
		r.db.Where("stage_id = ?", sid).Delete(&models.ReleaseStory{})
		r.db.Where("stage_id = ?", sid).Delete(&models.ReleaseSlackUpdate{})
	}
	r.db.Where("release_id = ?", id).Delete(&models.ReleaseStage{})
	return r.db.Delete(&models.Release{}, id).Error
}

func (r *releaseRepo) CreateStage(s models.ReleaseStage) error { return r.db.Create(&s).Error }

func (r *releaseRepo) DeleteStage(id uint) error {
	r.db.Where("stage_id = ?", id).Delete(&models.ReleaseStory{})
	r.db.Where("stage_id = ?", id).Delete(&models.ReleaseSlackUpdate{})
	return r.db.Delete(&models.ReleaseStage{}, id).Error
}

func (r *releaseRepo) UpdateStageStatus(id uint, status string) error {
	return r.db.Model(&models.ReleaseStage{}).Where("id = ?", id).Update("status", status).Error
}

func (r *releaseRepo) CreateStory(s models.ReleaseStory) error { return r.db.Create(&s).Error }

func (r *releaseRepo) DeleteStory(id uint) error {
	return r.db.Delete(&models.ReleaseStory{}, id).Error
}

func (r *releaseRepo) UpdateStoryStatus(id uint, status string) error {
	return r.db.Model(&models.ReleaseStory{}).Where("id = ?", id).Update("status", status).Error
}

func (r *releaseRepo) UpdateStory(id uint, title, assignee string) error {
	return r.db.Model(&models.ReleaseStory{}).Where("id = ?", id).
		Updates(map[string]interface{}{"title": title, "assignee": assignee}).Error
}

func (r *releaseRepo) CreateSlackUpdate(u models.ReleaseSlackUpdate) error {
	return r.db.Create(&u).Error
}

func (r *releaseRepo) DeleteSlackUpdate(id uint) error {
	return r.db.Delete(&models.ReleaseSlackUpdate{}, id).Error
}
