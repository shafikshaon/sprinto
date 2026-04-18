// Package service contains business logic that orchestrates repository calls.
// Handlers call services; services call repositories.
package service

import (
	"math"
	"strings"
	"time"

	"sprinto/models"
	"sprinto/repository"
)

// ─── Sprint ───────────────────────────────────────────────────────────────────

type SprintService interface {
	ActiveSprint() (models.Sprint, error)
	AddTask(sprintID uint, title, assignee, status, priority string) error
	RemoveTask(id uint) error
	UpdateProgress(sprintID uint, progress int) error
}

type sprintService struct{ repo repository.SprintRepository }

func NewSprintService(r repository.SprintRepository) SprintService {
	return &sprintService{repo: r}
}

func (s *sprintService) ActiveSprint() (models.Sprint, error) {
	return s.repo.ActiveSprint()
}

func (s *sprintService) AddTask(sprintID uint, title, assignee, status, priority string) error {
	if strings.TrimSpace(title) == "" {
		return nil
	}
	return s.repo.CreateTask(models.SprintTask{
		SprintID: sprintID,
		Title:    strings.TrimSpace(title),
		Assignee: strings.TrimSpace(assignee),
		Status:   status,
		Priority: priority,
	})
}

func (s *sprintService) RemoveTask(id uint) error { return s.repo.DeleteTask(id) }
func (s *sprintService) UpdateProgress(id uint, p int) error {
	return s.repo.UpdateProgress(id, clamp(p, 0, 100))
}

// ─── Standup ──────────────────────────────────────────────────────────────────

type DateNav struct {
	Raw     string
	Display string
}

type StandupService interface {
	ByDate(date string) ([]models.StandupEntry, error)
	Add(member, role, yesterday, today, blockers, status string) error
	Remove(id uint) error
	RecentDates(limit int) ([]DateNav, error)
}

type standupService struct{ repo repository.StandupRepository }

func NewStandupService(r repository.StandupRepository) StandupService {
	return &standupService{repo: r}
}

func (s *standupService) ByDate(date string) ([]models.StandupEntry, error) {
	return s.repo.ByDate(date)
}

func (s *standupService) Add(member, role, yesterday, today, blockers, status string) error {
	if strings.TrimSpace(member) == "" {
		return nil
	}
	if strings.TrimSpace(blockers) == "" {
		blockers = "None"
	}
	return s.repo.Create(models.StandupEntry{
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

func (s *standupService) RecentDates(limit int) ([]DateNav, error) {
	raw, err := s.repo.RecentDates(limit)
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
	All() ([]models.Deadline, error)
	Add(title, project, dueDateRaw, priority string) error
	Remove(id uint) error
}

type deadlineService struct{ repo repository.DeadlineRepository }

func NewDeadlineService(r repository.DeadlineRepository) DeadlineService {
	return &deadlineService{repo: r}
}

func (s *deadlineService) All() ([]models.Deadline, error) {
	deadlines, err := s.repo.All()
	if err != nil {
		return nil, err
	}
	for i := range deadlines {
		deadlines[i].DaysLeft = daysLeft(deadlines[i].DueDateRaw)
		deadlines[i].DueDate = formatDate(deadlines[i].DueDateRaw)
	}
	return deadlines, nil
}

func (s *deadlineService) Add(title, project, dueDateRaw, priority string) error {
	if strings.TrimSpace(title) == "" || dueDateRaw == "" {
		return nil
	}
	return s.repo.Create(models.Deadline{
		Title:      strings.TrimSpace(title),
		Project:    strings.TrimSpace(project),
		DueDateRaw: dueDateRaw,
		Priority:   priority,
	})
}

func (s *deadlineService) Remove(id uint) error { return s.repo.Delete(id) }

// ─── Meeting ──────────────────────────────────────────────────────────────────

type MeetingService interface {
	All() ([]models.Meeting, error)
	Add(title, date, attendeesCSV, notes string) error
	Remove(id uint) error
}

type meetingService struct{ repo repository.MeetingRepository }

func NewMeetingService(r repository.MeetingRepository) MeetingService {
	return &meetingService{repo: r}
}

func (s *meetingService) All() ([]models.Meeting, error) {
	meetings, err := s.repo.All()
	if err != nil {
		return nil, err
	}
	// Populate computed Attendees slice from CSV string
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

func (s *meetingService) Add(title, date, attendeesCSV, notes string) error {
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
		Title:       strings.TrimSpace(title),
		Date:        strings.TrimSpace(date),
		AttendeeCSV: strings.Join(attendees, ","),
		Notes:       strings.TrimSpace(notes),
	})
}

func (s *meetingService) Remove(id uint) error { return s.repo.Delete(id) }

// ─── Dev Task ─────────────────────────────────────────────────────────────────

type DevTaskService interface {
	All() ([]models.DevTask, error)
	Add(title, typ, assignee, status, priority string) error
	Remove(id uint) error
	OpenCountsByType() (map[string]int, error)
}

type devTaskService struct{ repo repository.DevTaskRepository }

func NewDevTaskService(r repository.DevTaskRepository) DevTaskService {
	return &devTaskService{repo: r}
}

func (s *devTaskService) All() ([]models.DevTask, error) { return s.repo.All() }

func (s *devTaskService) Add(title, typ, assignee, status, priority string) error {
	if strings.TrimSpace(title) == "" {
		return nil
	}
	return s.repo.Create(models.DevTask{
		Title:    strings.TrimSpace(title),
		Type:     typ,
		Assignee: strings.TrimSpace(assignee),
		Status:   status,
		Priority: priority,
	})
}

func (s *devTaskService) Remove(id uint) error { return s.repo.Delete(id) }
func (s *devTaskService) OpenCountsByType() (map[string]int, error) {
	return s.repo.OpenCountsByType()
}

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

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
