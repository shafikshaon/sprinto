// Package repository handles all direct database operations via GORM.
// No business logic lives here — only CRUD.
package repository

import (
	"gorm.io/gorm"

	"sprinto/models"
)

// ─── Sprint ───────────────────────────────────────────────────────────────────

type SprintRepository interface {
	ActiveSprint() (models.Sprint, error)
	CreateTask(t models.SprintTask) error
	DeleteTask(id uint) error
	UpdateProgress(sprintID uint, progress int) error
}

type sprintRepo struct{ db *gorm.DB }

func NewSprintRepository(db *gorm.DB) SprintRepository { return &sprintRepo{db: db} }

func (r *sprintRepo) ActiveSprint() (models.Sprint, error) {
	var s models.Sprint
	result := r.db.Preload("Tasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).Where("active = ?", true).Order("id DESC").First(&s)
	return s, result.Error
}

func (r *sprintRepo) CreateTask(t models.SprintTask) error {
	return r.db.Create(&t).Error
}

func (r *sprintRepo) DeleteTask(id uint) error {
	return r.db.Delete(&models.SprintTask{}, id).Error
}

func (r *sprintRepo) UpdateProgress(sprintID uint, progress int) error {
	return r.db.Model(&models.Sprint{}).Where("id = ?", sprintID).Update("progress", progress).Error
}

// ─── Standup ──────────────────────────────────────────────────────────────────

type StandupRepository interface {
	ByDate(date string) ([]models.StandupEntry, error)
	Create(e models.StandupEntry) error
	Delete(id uint) error
	RecentDates(limit int) ([]string, error)
}

type standupRepo struct{ db *gorm.DB }

func NewStandupRepository(db *gorm.DB) StandupRepository { return &standupRepo{db: db} }

func (r *standupRepo) ByDate(date string) ([]models.StandupEntry, error) {
	var entries []models.StandupEntry
	result := r.db.Where("date = ?", date).Order("created_at ASC").Find(&entries)
	return entries, result.Error
}

func (r *standupRepo) Create(e models.StandupEntry) error {
	return r.db.Create(&e).Error
}

func (r *standupRepo) Delete(id uint) error {
	return r.db.Delete(&models.StandupEntry{}, id).Error
}

func (r *standupRepo) RecentDates(limit int) ([]string, error) {
	var dates []string
	result := r.db.Model(&models.StandupEntry{}).
		Select("DISTINCT date").
		Order("date DESC").
		Limit(limit).
		Pluck("date", &dates)
	return dates, result.Error
}

// ─── Deadline ─────────────────────────────────────────────────────────────────

type DeadlineRepository interface {
	All() ([]models.Deadline, error)
	Create(d models.Deadline) error
	Delete(id uint) error
}

type deadlineRepo struct{ db *gorm.DB }

func NewDeadlineRepository(db *gorm.DB) DeadlineRepository { return &deadlineRepo{db: db} }

func (r *deadlineRepo) All() ([]models.Deadline, error) {
	var dl []models.Deadline
	result := r.db.Order("due_date ASC").Find(&dl)
	return dl, result.Error
}

func (r *deadlineRepo) Create(d models.Deadline) error {
	return r.db.Create(&d).Error
}

func (r *deadlineRepo) Delete(id uint) error {
	return r.db.Delete(&models.Deadline{}, id).Error
}

// ─── Meeting ──────────────────────────────────────────────────────────────────

type MeetingRepository interface {
	All() ([]models.Meeting, error)
	Create(m models.Meeting) error
	Delete(id uint) error
}

type meetingRepo struct{ db *gorm.DB }

func NewMeetingRepository(db *gorm.DB) MeetingRepository { return &meetingRepo{db: db} }

func (r *meetingRepo) All() ([]models.Meeting, error) {
	var meetings []models.Meeting
	result := r.db.Preload("ActionItems").Order("created_at DESC").Find(&meetings)
	return meetings, result.Error
}

func (r *meetingRepo) Create(m models.Meeting) error {
	return r.db.Create(&m).Error
}

func (r *meetingRepo) Delete(id uint) error {
	return r.db.Delete(&models.Meeting{}, id).Error
}

// ─── Dev Task ─────────────────────────────────────────────────────────────────

type DevTaskRepository interface {
	All() ([]models.DevTask, error)
	Create(t models.DevTask) error
	Delete(id uint) error
	OpenCountsByType() (map[string]int, error)
}

type devTaskRepo struct{ db *gorm.DB }

func NewDevTaskRepository(db *gorm.DB) DevTaskRepository { return &devTaskRepo{db: db} }

func (r *devTaskRepo) All() ([]models.DevTask, error) {
	var tasks []models.DevTask
	result := r.db.Order("created_at DESC").Find(&tasks)
	return tasks, result.Error
}

func (r *devTaskRepo) Create(t models.DevTask) error {
	return r.db.Create(&t).Error
}

func (r *devTaskRepo) Delete(id uint) error {
	return r.db.Delete(&models.DevTask{}, id).Error
}

func (r *devTaskRepo) OpenCountsByType() (map[string]int, error) {
	type row struct {
		Type  string
		Count int
	}
	var rows []row
	result := r.db.Model(&models.DevTask{}).
		Select("type, COUNT(*) as count").
		Where("status != ?", "Done").
		Group("type").
		Scan(&rows)
	if result.Error != nil {
		return nil, result.Error
	}
	counts := map[string]int{"Improvement": 0, "Tech Debt": 0, "Research": 0}
	for _, r := range rows {
		counts[r.Type] = r.Count
	}
	return counts, nil
}
