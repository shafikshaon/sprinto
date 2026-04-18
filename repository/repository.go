// Package repository handles all direct database operations via GORM.
// No business logic lives here — only CRUD.
package repository

import (
	"gorm.io/gorm"

	"sprinto/models"
)

// ─── Sprint ───────────────────────────────────────────────────────────────────

type SprintRepository interface {
	ActiveSprint(projectID uint) (models.Sprint, error)
	CreateTask(t models.SprintTask) error
	DeleteTask(id uint) error
	UpdateProgress(sprintID uint, progress int) error
}

type sprintRepo struct{ db *gorm.DB }

func NewSprintRepository(db *gorm.DB) SprintRepository { return &sprintRepo{db: db} }

func (r *sprintRepo) ActiveSprint(projectID uint) (models.Sprint, error) {
	var s models.Sprint
	q := r.db.Preload("Tasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).Where("active = ?", true)
	if projectID > 0 {
		q = q.Where("project_id = ?", projectID)
	}
	result := q.Order("id DESC").First(&s)
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
	ByDate(date string, projectID uint) ([]models.StandupEntry, error)
	Create(e models.StandupEntry) error
	Delete(id uint) error
	RecentDates(limit int, projectID uint) ([]string, error)
}

type standupRepo struct{ db *gorm.DB }

func NewStandupRepository(db *gorm.DB) StandupRepository { return &standupRepo{db: db} }

func (r *standupRepo) ByDate(date string, projectID uint) ([]models.StandupEntry, error) {
	var entries []models.StandupEntry
	q := r.db.Where("date = ?", date)
	if projectID > 0 {
		q = q.Where("project_id = ?", projectID)
	}
	result := q.Order("created_at ASC").Find(&entries)
	return entries, result.Error
}

func (r *standupRepo) Create(e models.StandupEntry) error {
	return r.db.Create(&e).Error
}

func (r *standupRepo) Delete(id uint) error {
	return r.db.Delete(&models.StandupEntry{}, id).Error
}

func (r *standupRepo) RecentDates(limit int, projectID uint) ([]string, error) {
	var dates []string
	q := r.db.Model(&models.StandupEntry{}).Select("DISTINCT date")
	if projectID > 0 {
		q = q.Where("project_id = ?", projectID)
	}
	result := q.Order("date DESC").Limit(limit).Pluck("date", &dates)
	return dates, result.Error
}

// ─── Deadline ─────────────────────────────────────────────────────────────────

type DeadlineRepository interface {
	All(projectID uint) ([]models.Deadline, error)
	Create(d models.Deadline) error
	Delete(id uint) error
}

type deadlineRepo struct{ db *gorm.DB }

func NewDeadlineRepository(db *gorm.DB) DeadlineRepository { return &deadlineRepo{db: db} }

func (r *deadlineRepo) All(projectID uint) ([]models.Deadline, error) {
	var dl []models.Deadline
	q := r.db.Order("due_date ASC")
	if projectID > 0 {
		q = q.Where("project_id = ?", projectID)
	}
	result := q.Find(&dl)
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
		// Show meetings for this project AND global meetings (project_id = 0)
		q = q.Where("project_id = ? OR project_id = 0", projectID)
	}
	result := q.Find(&meetings)
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
	All(projectID uint) ([]models.DevTask, error)
	Create(t models.DevTask) error
	Delete(id uint) error
	OpenCountsByType(projectID uint) (map[string]int, error)
}

type devTaskRepo struct{ db *gorm.DB }

func NewDevTaskRepository(db *gorm.DB) DevTaskRepository { return &devTaskRepo{db: db} }

func (r *devTaskRepo) All(projectID uint) ([]models.DevTask, error) {
	var tasks []models.DevTask
	q := r.db.Order("created_at DESC")
	if projectID > 0 {
		q = q.Where("project_id = ?", projectID)
	}
	result := q.Find(&tasks)
	return tasks, result.Error
}

func (r *devTaskRepo) Create(t models.DevTask) error {
	return r.db.Create(&t).Error
}

func (r *devTaskRepo) Delete(id uint) error {
	return r.db.Delete(&models.DevTask{}, id).Error
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
	result := q.Group("type").Scan(&rows)
	if result.Error != nil {
		return nil, result.Error
	}
	counts := map[string]int{"Improvement": 0, "Tech Debt": 0, "Research": 0}
	for _, r := range rows {
		counts[r.Type] = r.Count
	}
	return counts, nil
}

// ─── Release ──────────────────────────────────────────────────────────────────

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
	result := q.Find(&releases)
	return releases, result.Error
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

func (r *releaseRepo) CreateSlackUpdate(u models.ReleaseSlackUpdate) error {
	return r.db.Create(&u).Error
}

func (r *releaseRepo) DeleteSlackUpdate(id uint) error {
	return r.db.Delete(&models.ReleaseSlackUpdate{}, id).Error
}

// ─── Project ──────────────────────────────────────────────────────────────────

type ProjectRepository interface {
	All() ([]models.Project, error)
	AllWithMembers() ([]models.Project, error)
	ByID(id uint) (models.Project, error)
	Create(p models.Project) error
	Delete(id uint) error
	AddMember(projectID, memberID uint) error
	RemoveMember(projectID, memberID uint) error
}

type projectRepo struct{ db *gorm.DB }

func NewProjectRepository(db *gorm.DB) ProjectRepository { return &projectRepo{db: db} }

func (r *projectRepo) All() ([]models.Project, error) {
	var projects []models.Project
	result := r.db.Order("created_at ASC").Find(&projects)
	return projects, result.Error
}

func (r *projectRepo) AllWithMembers() ([]models.Project, error) {
	var projects []models.Project
	result := r.db.Preload("Members").Order("created_at ASC").Find(&projects)
	return projects, result.Error
}

func (r *projectRepo) ByID(id uint) (models.Project, error) {
	var p models.Project
	result := r.db.Preload("Members").First(&p, id)
	return p, result.Error
}

func (r *projectRepo) Create(p models.Project) error { return r.db.Create(&p).Error }

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

// ─── Team Member ──────────────────────────────────────────────────────────────

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
	result := r.db.Preload("Projects").Order("created_at ASC").Find(&members)
	return members, result.Error
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
