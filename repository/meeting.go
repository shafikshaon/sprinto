package repository

import (
	"gorm.io/gorm"
	"sprinto/models"
)

type MeetingFilter struct {
	Search   string
	DateFrom string
	DateTo   string
	Project  uint
}

type MeetingRepository interface {
	All(f MeetingFilter, page, perPage int) ([]models.Meeting, int64, error)
	Create(m models.Meeting) error
	Update(id uint, projectID uint, title, date, attendees, notes string) error
	Delete(id uint) error
}

type meetingRepo struct{ db *gorm.DB }

func NewMeetingRepository(db *gorm.DB) MeetingRepository { return &meetingRepo{db: db} }

func (r *meetingRepo) filtered(f MeetingFilter) *gorm.DB {
	q := r.db.Model(&models.Meeting{})
	if f.Project > 0 {
		q = q.Where("project_id = ?", f.Project)
	}
	if f.Search != "" {
		q = q.Where("title LIKE ? OR notes LIKE ?", "%"+f.Search+"%", "%"+f.Search+"%")
	}
	if f.DateFrom != "" {
		q = q.Where("date >= ?", f.DateFrom)
	}
	if f.DateTo != "" {
		q = q.Where("date <= ?", f.DateTo)
	}
	return q
}

func (r *meetingRepo) All(f MeetingFilter, page, perPage int) ([]models.Meeting, int64, error) {
	var meetings []models.Meeting
	var total int64
	if err := r.filtered(f).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	err := r.filtered(f).
		Preload("ActionItems").Preload("Project").
		Order("created_at DESC").
		Offset(offset).Limit(perPage).
		Find(&meetings).Error
	return meetings, total, err
}

func (r *meetingRepo) Create(m models.Meeting) error { return r.db.Create(&m).Error }

func (r *meetingRepo) Update(id uint, projectID uint, title, date, attendees, notes string) error {
	return r.db.Model(&models.Meeting{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"project_id": projectID,
			"title":      title,
			"date":       date,
			"attendees":  attendees,
			"notes":      notes,
		}).Error
}

func (r *meetingRepo) Delete(id uint) error { return r.db.Delete(&models.Meeting{}, id).Error }
