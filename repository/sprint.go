package repository

import (
	"gorm.io/gorm"
	"sprinto/models"
)

type SprintRepository interface {
	// Sprint management
	ListSprints(projectID uint) ([]models.Sprint, error)
	CreateSprint(s models.Sprint) error
	UpdateSprint(id uint, name, goal, startDate, endDate string) error
	DeleteSprint(id uint) error
	ActivateSprint(id uint, projectID uint) error
	// Board
	ActiveSprint(projectID uint) (models.Sprint, error)
	TaskByID(id uint) (models.Task, error)
	CreateTask(t models.Task, assigneeIDs []uint) error
	UpdateTask(id uint, title, status, priority string, assigneeIDs []uint) error
	DeleteTask(id uint) error
	UpdateProgress(sprintID uint, progress int) error
	AddComment(c models.TaskComment) error
	DeleteComment(id uint) error
	// Release (merged)
	UpdateRelease(id uint, description, status, targetDate string) error
	CreateStage(s models.ReleaseStage) error
	DeleteStage(id uint) error
	UpdateStageStatus(id uint, status string) error
	CreateStory(t models.Task) error
	DeleteStory(id uint) error
	UpdateStoryStatus(id uint, status string) error
	UpdateStory(id uint, title string, assigneeID *uint) error
	CreateSlackUpdate(u models.ReleaseSlackUpdate) error
	DeleteSlackUpdate(id uint) error
}

type sprintRepo struct{ db *gorm.DB }

func NewSprintRepository(db *gorm.DB) SprintRepository { return &sprintRepo{db: db} }

func (r *sprintRepo) ListSprints(projectID uint) ([]models.Sprint, error) {
	var sprints []models.Sprint
	q := r.db.Order("created_at DESC")
	if projectID > 0 {
		q = q.Where("project_id = ?", projectID)
	}
	return sprints, q.Find(&sprints).Error
}

func (r *sprintRepo) CreateSprint(s models.Sprint) error { return r.db.Create(&s).Error }

func (r *sprintRepo) UpdateSprint(id uint, name, goal, startDate, endDate string) error {
	return r.db.Model(&models.Sprint{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":       name,
			"goal":       goal,
			"start_date": startDate,
			"end_date":   endDate,
		}).Error
}

func (r *sprintRepo) DeleteSprint(id uint) error {
	r.db.Where("sprint_id = ?", id).Delete(&models.Task{})
	r.db.Where("sprint_id = ?", id).Delete(&models.ReleaseStage{})
	return r.db.Delete(&models.Sprint{}, id).Error
}

func (r *sprintRepo) ActivateSprint(id uint, projectID uint) error {
	if err := r.db.Model(&models.Sprint{}).
		Where("project_id = ?", projectID).
		Update("active", false).Error; err != nil {
		return err
	}
	return r.db.Model(&models.Sprint{}).Where("id = ?", id).Update("active", true).Error
}

func (r *sprintRepo) ActiveSprint(projectID uint) (models.Sprint, error) {
	var s models.Sprint
	q := r.db.
		Preload("Tasks", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
		Preload("Tasks.Assignees").
		Preload("Tasks.Comments").
		Preload("Stages", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
		Preload("Stages.Stories", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
		Preload("Stages.Stories.Assignee").
		Preload("Stages.SlackUpdates", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
		Where("active = ?", true)
	if projectID > 0 {
		q = q.Where("project_id = ?", projectID)
	}
	result := q.Order("id DESC").First(&s)
	return s, result.Error
}

func (r *sprintRepo) TaskByID(id uint) (models.Task, error) {
	var t models.Task
	result := r.db.Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).Preload("Assignees").First(&t, id)
	return t, result.Error
}

func (r *sprintRepo) CreateTask(t models.Task, assigneeIDs []uint) error {
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
	if err := r.db.Model(&models.Task{}).Where("id = ?", id).
		Updates(map[string]interface{}{"title": title, "status": status, "priority": priority}).Error; err != nil {
		return err
	}
	task := models.Task{Model: gorm.Model{ID: id}}
	members := make([]models.TeamMember, len(assigneeIDs))
	for i, aid := range assigneeIDs {
		members[i] = models.TeamMember{Model: gorm.Model{ID: aid}}
	}
	return r.db.Model(&task).Association("Assignees").Replace(members)
}

func (r *sprintRepo) DeleteTask(id uint) error { return r.db.Delete(&models.Task{}, id).Error }

func (r *sprintRepo) UpdateProgress(sprintID uint, progress int) error {
	return r.db.Model(&models.Sprint{}).Where("id = ?", sprintID).Update("progress", progress).Error
}

func (r *sprintRepo) AddComment(c models.TaskComment) error { return r.db.Create(&c).Error }

func (r *sprintRepo) DeleteComment(id uint) error {
	return r.db.Delete(&models.TaskComment{}, id).Error
}

// ── Release methods ───────────────────────────────────────────────────────────

func (r *sprintRepo) UpdateRelease(id uint, description, status, targetDate string) error {
	return r.db.Model(&models.Sprint{}).Where("id = ?", id).
		Updates(map[string]interface{}{"description": description, "status": status, "target_date": targetDate}).Error
}

func (r *sprintRepo) CreateStage(s models.ReleaseStage) error { return r.db.Create(&s).Error }

func (r *sprintRepo) DeleteStage(id uint) error {
	r.db.Where("release_stage_id = ?", id).Delete(&models.Task{})
	r.db.Where("stage_id = ?", id).Delete(&models.ReleaseSlackUpdate{})
	return r.db.Delete(&models.ReleaseStage{}, id).Error
}

func (r *sprintRepo) UpdateStageStatus(id uint, status string) error {
	return r.db.Model(&models.ReleaseStage{}).Where("id = ?", id).Update("status", status).Error
}

func (r *sprintRepo) CreateStory(t models.Task) error { return r.db.Create(&t).Error }

func (r *sprintRepo) DeleteStory(id uint) error { return r.db.Delete(&models.Task{}, id).Error }

func (r *sprintRepo) UpdateStoryStatus(id uint, status string) error {
	return r.db.Model(&models.Task{}).Where("id = ?", id).Update("status", status).Error
}

func (r *sprintRepo) UpdateStory(id uint, title string, assigneeID *uint) error {
	return r.db.Model(&models.Task{}).Where("id = ?", id).
		Updates(map[string]interface{}{"title": title, "assignee_id": assigneeID}).Error
}

func (r *sprintRepo) CreateSlackUpdate(u models.ReleaseSlackUpdate) error {
	return r.db.Create(&u).Error
}

func (r *sprintRepo) DeleteSlackUpdate(id uint) error {
	return r.db.Delete(&models.ReleaseSlackUpdate{}, id).Error
}
