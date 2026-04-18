// Package service contains business logic that orchestrates repository calls.
// Handlers call services; services call repositories.
package service

import (
	"errors"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"

	"sprinto/models"
	"sprinto/repository"
)

// ─── Sprint ───────────────────────────────────────────────────────────────────

type SprintService interface {
	ActiveSprint(projectID uint) (models.Sprint, error)
	TaskByID(id uint) (models.SprintTask, error)
	AddTask(sprintID uint, title string, assignees []string, status, priority string) error
	RemoveTask(id uint) error
	UpdateProgress(sprintID uint, progress int) error
	AddComment(taskID uint, author, content string) error
	DeleteComment(id uint) error
}

type sprintService struct{ repo repository.SprintRepository }

func NewSprintService(r repository.SprintRepository) SprintService {
	return &sprintService{repo: r}
}

func (s *sprintService) ActiveSprint(projectID uint) (models.Sprint, error) {
	sprint, err := s.repo.ActiveSprint(projectID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Sprint{}, nil
	}
	if err == nil {
		for i := range sprint.Tasks {
			sprint.Tasks[i].Assignees = splitAssignees(sprint.Tasks[i].AssigneeCSV)
		}
	}
	return sprint, err
}

func (s *sprintService) AddTask(sprintID uint, title string, assignees []string, status, priority string) error {
	if strings.TrimSpace(title) == "" {
		return nil
	}
	return s.repo.CreateTask(models.SprintTask{
		SprintID:    sprintID,
		Title:       strings.TrimSpace(title),
		AssigneeCSV: strings.Join(assignees, ","),
		Status:      status,
		Priority:    priority,
	})
}

func (s *sprintService) TaskByID(id uint) (models.SprintTask, error) {
	task, err := s.repo.TaskByID(id)
	if err == nil {
		task.Assignees = splitAssignees(task.AssigneeCSV)
	}
	return task, err
}

func splitAssignees(csv string) []string {
	var out []string
	for _, p := range strings.Split(csv, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func (s *sprintService) RemoveTask(id uint) error { return s.repo.DeleteTask(id) }

func (s *sprintService) UpdateProgress(id uint, p int) error {
	return s.repo.UpdateProgress(id, clamp(p, 0, 100))
}

func (s *sprintService) AddComment(taskID uint, author, content string) error {
	author = strings.TrimSpace(author)
	content = strings.TrimSpace(content)
	if content == "" {
		return nil
	}
	if author == "" {
		author = "Anonymous"
	}
	return s.repo.AddComment(models.SprintTaskComment{
		TaskID:  taskID,
		Author:  author,
		Content: content,
	})
}

func (s *sprintService) DeleteComment(id uint) error { return s.repo.DeleteComment(id) }

// ─── Standup ──────────────────────────────────────────────────────────────────

type DateNav struct {
	Raw     string
	Display string
}

type StandupService interface {
	ByDate(date string, projectID uint) ([]models.StandupEntry, error)
	Add(member, role, yesterday, today, blockers, status string, projectID uint) error
	Remove(id uint) error
	RecentDates(limit int, projectID uint) ([]DateNav, error)
}

type standupService struct{ repo repository.StandupRepository }

func NewStandupService(r repository.StandupRepository) StandupService {
	return &standupService{repo: r}
}

func (s *standupService) ByDate(date string, projectID uint) ([]models.StandupEntry, error) {
	return s.repo.ByDate(date, projectID)
}

func (s *standupService) Add(member, role, yesterday, today, blockers, status string, projectID uint) error {
	if strings.TrimSpace(member) == "" {
		return nil
	}
	if strings.TrimSpace(blockers) == "" {
		blockers = "None"
	}
	return s.repo.Create(models.StandupEntry{
		ProjectID: projectID,
		Member:    strings.TrimSpace(member),
		Role:      strings.TrimSpace(role),
		Yesterday: strings.TrimSpace(yesterday),
		Today:     strings.TrimSpace(today),
		Blockers:  strings.TrimSpace(blockers),
		Status:    status,
		Date:      time.Now().Format("2006-01-02"),
	})
}

func (s *standupService) Remove(id uint) error { return s.repo.Delete(id) }

func (s *standupService) RecentDates(limit int, projectID uint) ([]DateNav, error) {
	raw, err := s.repo.RecentDates(limit, projectID)
	if err != nil {
		return nil, err
	}
	nav := make([]DateNav, 0, len(raw))
	for _, d := range raw {
		t, _ := time.Parse("2006-01-02", d)
		nav = append(nav, DateNav{Raw: d, Display: t.Format("Jan 2")})
	}
	return nav, nil
}

// ─── Deadline ─────────────────────────────────────────────────────────────────

type DeadlineService interface {
	All(projectID uint) ([]models.Deadline, error)
	Add(title, project, dueDateRaw, priority string, projectID uint) error
	Remove(id uint) error
}

type deadlineService struct{ repo repository.DeadlineRepository }

func NewDeadlineService(r repository.DeadlineRepository) DeadlineService {
	return &deadlineService{repo: r}
}

func (s *deadlineService) All(projectID uint) ([]models.Deadline, error) {
	deadlines, err := s.repo.All(projectID)
	if err != nil {
		return nil, err
	}
	for i := range deadlines {
		deadlines[i].DaysLeft = daysLeft(deadlines[i].DueDateRaw)
		deadlines[i].DueDate = formatDate(deadlines[i].DueDateRaw)
	}
	return deadlines, nil
}

func (s *deadlineService) Add(title, project, dueDateRaw, priority string, projectID uint) error {
	if strings.TrimSpace(title) == "" || dueDateRaw == "" {
		return nil
	}
	return s.repo.Create(models.Deadline{
		ProjectID:  projectID,
		Title:      strings.TrimSpace(title),
		Project:    strings.TrimSpace(project),
		DueDateRaw: dueDateRaw,
		Priority:   priority,
	})
}

func (s *deadlineService) Remove(id uint) error { return s.repo.Delete(id) }

// ─── Meeting ──────────────────────────────────────────────────────────────────

type MeetingService interface {
	All(projectID uint) ([]models.Meeting, error)
	Add(title, date, attendeesCSV, notes string, projectID uint) error
	Remove(id uint) error
}

type meetingService struct{ repo repository.MeetingRepository }

func NewMeetingService(r repository.MeetingRepository) MeetingService {
	return &meetingService{repo: r}
}

func (s *meetingService) All(projectID uint) ([]models.Meeting, error) {
	meetings, err := s.repo.All(projectID)
	if err != nil {
		return nil, err
	}
	for i := range meetings {
		if meetings[i].AttendeeCSV != "" {
			for _, a := range strings.Split(meetings[i].AttendeeCSV, ",") {
				if t := strings.TrimSpace(a); t != "" {
					meetings[i].Attendees = append(meetings[i].Attendees, t)
				}
			}
		}
	}
	return meetings, nil
}

func (s *meetingService) Add(title, date, attendeesCSV, notes string, projectID uint) error {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(date) == "" {
		return nil
	}
	var attendees []string
	for _, a := range strings.Split(attendeesCSV, ",") {
		if t := strings.TrimSpace(a); t != "" {
			attendees = append(attendees, t)
		}
	}
	return s.repo.Create(models.Meeting{
		ProjectID:   projectID,
		Title:       strings.TrimSpace(title),
		Date:        strings.TrimSpace(date),
		AttendeeCSV: strings.Join(attendees, ","),
		Notes:       strings.TrimSpace(notes),
	})
}

func (s *meetingService) Remove(id uint) error { return s.repo.Delete(id) }

// ─── Dev Task ─────────────────────────────────────────────────────────────────

type DevTaskService interface {
	All(projectID uint) ([]models.DevTask, error)
	ByID(id uint) (models.DevTask, error)
	Add(title, typ string, assignees []string, status, priority string, projectID uint) error
	Remove(id uint) error
	OpenCountsByType(projectID uint) (map[string]int, error)
	AddComment(taskID uint, author, content string) error
	DeleteComment(id uint) error
}

type devTaskService struct{ repo repository.DevTaskRepository }

func NewDevTaskService(r repository.DevTaskRepository) DevTaskService {
	return &devTaskService{repo: r}
}

func (s *devTaskService) All(projectID uint) ([]models.DevTask, error) {
	tasks, err := s.repo.All(projectID)
	if err == nil {
		for i := range tasks {
			tasks[i].Assignees = splitAssignees(tasks[i].AssigneeCSV)
		}
	}
	return tasks, err
}

func (s *devTaskService) Add(title, typ string, assignees []string, status, priority string, projectID uint) error {
	if strings.TrimSpace(title) == "" {
		return nil
	}
	return s.repo.Create(models.DevTask{
		ProjectID:   projectID,
		Title:       strings.TrimSpace(title),
		Type:        typ,
		AssigneeCSV: strings.Join(assignees, ","),
		Status:      status,
		Priority:    priority,
	})
}

func (s *devTaskService) ByID(id uint) (models.DevTask, error) {
	task, err := s.repo.ByID(id)
	if err == nil {
		task.Assignees = splitAssignees(task.AssigneeCSV)
	}
	return task, err
}

func (s *devTaskService) Remove(id uint) error { return s.repo.Delete(id) }

func (s *devTaskService) OpenCountsByType(projectID uint) (map[string]int, error) {
	return s.repo.OpenCountsByType(projectID)
}

func (s *devTaskService) AddComment(taskID uint, author, content string) error {
	author = strings.TrimSpace(author)
	content = strings.TrimSpace(content)
	if content == "" {
		return nil
	}
	if author == "" {
		author = "Anonymous"
	}
	return s.repo.AddComment(models.DevTaskComment{
		TaskID:  taskID,
		Author:  author,
		Content: content,
	})
}

func (s *devTaskService) DeleteComment(id uint) error { return s.repo.DeleteComment(id) }

// ─── Release ──────────────────────────────────────────────────────────────────

type ReleaseService interface {
	All(projectID uint) ([]models.Release, error)
	ByID(id uint) (models.Release, error)
	Create(name, description, status, targetDate string, projectID uint) error
	Delete(id uint) error
	AddStage(releaseID uint, name, status string) error
	DeleteStage(id uint) error
	UpdateStageStatus(id uint, status string) error
	AddStory(stageID uint, title, assignee string) error
	DeleteStory(id uint) error
	UpdateStoryStatus(id uint, status string) error
	AddSlackUpdate(stageID uint, channel, message, author string) error
	DeleteSlackUpdate(id uint) error
}

type releaseService struct{ repo repository.ReleaseRepository }

func NewReleaseService(r repository.ReleaseRepository) ReleaseService {
	return &releaseService{repo: r}
}

func (s *releaseService) All(projectID uint) ([]models.Release, error) {
	return s.repo.All(projectID)
}
func (s *releaseService) ByID(id uint) (models.Release, error) { return s.repo.ByID(id) }

func (s *releaseService) Create(name, description, status, targetDate string, projectID uint) error {
	if strings.TrimSpace(name) == "" {
		return nil
	}
	return s.repo.Create(models.Release{
		ProjectID:   projectID,
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		Status:      status,
		TargetDate:  targetDate,
	})
}

func (s *releaseService) Delete(id uint) error { return s.repo.Delete(id) }

func (s *releaseService) AddStage(releaseID uint, name, status string) error {
	if strings.TrimSpace(name) == "" {
		return nil
	}
	return s.repo.CreateStage(models.ReleaseStage{
		ReleaseID: releaseID,
		Name:      strings.TrimSpace(name),
		Status:    status,
	})
}

func (s *releaseService) DeleteStage(id uint) error { return s.repo.DeleteStage(id) }
func (s *releaseService) UpdateStageStatus(id uint, status string) error {
	return s.repo.UpdateStageStatus(id, status)
}

func (s *releaseService) AddStory(stageID uint, title, assignee string) error {
	if strings.TrimSpace(title) == "" {
		return nil
	}
	return s.repo.CreateStory(models.ReleaseStory{
		StageID:  stageID,
		Title:    strings.TrimSpace(title),
		Assignee: strings.TrimSpace(assignee),
		Status:   "Pending",
	})
}

func (s *releaseService) DeleteStory(id uint) error { return s.repo.DeleteStory(id) }
func (s *releaseService) UpdateStoryStatus(id uint, status string) error {
	return s.repo.UpdateStoryStatus(id, status)
}

func (s *releaseService) AddSlackUpdate(stageID uint, channel, message, author string) error {
	if strings.TrimSpace(message) == "" {
		return nil
	}
	return s.repo.CreateSlackUpdate(models.ReleaseSlackUpdate{
		StageID:  stageID,
		Channel:  strings.TrimSpace(channel),
		Message:  strings.TrimSpace(message),
		Author:   strings.TrimSpace(author),
		PostedAt: time.Now().Format("Jan 2, 3:04 PM"),
	})
}

func (s *releaseService) DeleteSlackUpdate(id uint) error { return s.repo.DeleteSlackUpdate(id) }

// ─── Helpers ─────────────────────────────────────────────────────────────────

func daysLeft(raw string) int {
	due, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return 0
	}
	d := int(math.Ceil(time.Until(due).Hours() / 24))
	if d < 0 {
		return 0
	}
	return d
}

func formatDate(raw string) string {
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return raw
	}
	return t.Format("Jan 2, 2006")
}

// ─── Project ──────────────────────────────────────────────────────────────────

type ProjectService interface {
	All() ([]models.Project, error)
	AllWithMembers() ([]models.Project, error)
	Create(name, description string) error
	Delete(id uint) error
	AddMember(projectID, memberID uint) error
	RemoveMember(projectID, memberID uint) error
}

type projectService struct{ repo repository.ProjectRepository }

func NewProjectService(r repository.ProjectRepository) ProjectService {
	return &projectService{repo: r}
}

func (s *projectService) All() ([]models.Project, error) { return s.repo.All() }
func (s *projectService) AllWithMembers() ([]models.Project, error) {
	return s.repo.AllWithMembers()
}

func (s *projectService) Create(name, description string) error {
	if strings.TrimSpace(name) == "" {
		return nil
	}
	return s.repo.Create(models.Project{
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
	})
}

func (s *projectService) Delete(id uint) error         { return s.repo.Delete(id) }
func (s *projectService) AddMember(pid, mid uint) error { return s.repo.AddMember(pid, mid) }
func (s *projectService) RemoveMember(pid, mid uint) error {
	return s.repo.RemoveMember(pid, mid)
}

// ─── Team Member ──────────────────────────────────────────────────────────────

type TeamMemberService interface {
	All() ([]models.TeamMember, error)
	Create(name, role, email string) error
	Delete(id uint) error
}

type teamMemberService struct{ repo repository.TeamMemberRepository }

func NewTeamMemberService(r repository.TeamMemberRepository) TeamMemberService {
	return &teamMemberService{repo: r}
}

func (s *teamMemberService) All() ([]models.TeamMember, error) { return s.repo.All() }

func (s *teamMemberService) Create(name, role, email string) error {
	if strings.TrimSpace(name) == "" {
		return nil
	}
	return s.repo.Create(models.TeamMember{
		Name:  strings.TrimSpace(name),
		Role:  strings.TrimSpace(role),
		Email: strings.TrimSpace(email),
	})
}

func (s *teamMemberService) Delete(id uint) error { return s.repo.Delete(id) }

// ─────────────────────────────────────────────────────────────────────────────

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
