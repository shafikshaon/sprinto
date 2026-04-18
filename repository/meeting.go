package repository

import (
	"gorm.io/gorm"
	"sprinto/models"
)

type MeetingRepository interface {
	All(projectID uint) ([]models.Meeting, error)
	Create(m models.Meeting) error
	Delete(id uint) error
}

type meetingRepo struct{ db *gorm.DB }

func NewMeetingRepository(db *gorm.DB) MeetingRepository { return &meetingRepo{db: db} }

func (r *meetingRepo) All(projectID uint) ([]models.Meeting, error) {
	var meetings []models.Meeting
	q := r.db.Preload("ActionItems").Order("created_at DESC")
	if projectID > 0 {
		q = q.Where("project_id = ? OR project_id = 0", projectID)
	}
	return meetings, q.Find(&meetings).Error
}

func (r *meetingRepo) Create(m models.Meeting) error { return r.db.Create(&m).Error }

func (r *meetingRepo) Delete(id uint) error { return r.db.Delete(&models.Meeting{}, id).Error }
