package models

import "gorm.io/gorm"

// ─── Sprint ───────────────────────────────────────────────────────────────────

type Sprint struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Goal      string `gorm:"default:''"`
	Progress  int    `gorm:"default:0"`
	StartDate string
	EndDate   string
	Active    bool         `gorm:"default:false;index"`
	Tasks     []SprintTask `gorm:"foreignKey:SprintID"`
}

type SprintTask struct {
	gorm.Model
	SprintID uint   `gorm:"not null;index"`
	Title    string `gorm:"not null"`
	Assignee string
	Status   string `gorm:"default:'Todo'"`
	Priority string `gorm:"default:'Medium'"`
}

type SprintStats struct {
	Todo       int
	InProgress int
	Done       int
	Blocked    int
}

func ComputeStats(tasks []SprintTask) SprintStats {
	s := SprintStats{}
	for _, t := range tasks {
		switch t.Status {
		case "Todo":
			s.Todo++
		case "In Progress":
			s.InProgress++
		case "Done":
			s.Done++
		case "Blocked":
			s.Blocked++
		}
	}
	return s
}

// ─── Standup ──────────────────────────────────────────────────────────────────

type StandupEntry struct {
	gorm.Model
	Member    string `gorm:"not null"`
	Role      string
	Yesterday string
	Today     string
	Blockers  string `gorm:"default:'None'"`
	Status    string `gorm:"default:'On Track'"`
	Date      string `gorm:"not null;index"` // "YYYY-MM-DD"
}

func (StandupEntry) TableName() string { return "standups" }

// ─── Deadline ─────────────────────────────────────────────────────────────────

type Deadline struct {
	gorm.Model
	Title      string `gorm:"not null"`
	Project    string
	DueDateRaw string `gorm:"column:due_date;not null"` // "YYYY-MM-DD"
	Priority   string `gorm:"default:'Medium'"`
	// Computed by service — not persisted
	DueDate  string `gorm:"-"`
	DaysLeft int    `gorm:"-"`
}

// ─── Meeting ──────────────────────────────────────────────────────────────────

type ActionItem struct {
	gorm.Model
	MeetingID uint   `gorm:"not null;index"`
	Task      string `gorm:"not null"`
	Owner     string
	DueDate   string `gorm:"column:due_date"`
	Done      bool   `gorm:"default:false"`
}

func (ActionItem) TableName() string { return "meeting_action_items" }

type Meeting struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Date        string `gorm:"not null"`
	AttendeeCSV string `gorm:"column:attendees"`
	Notes       string
	ActionItems []ActionItem `gorm:"foreignKey:MeetingID"`
	// Computed by service — not persisted
	Attendees []string `gorm:"-"`
}

// ─── Dev Task ─────────────────────────────────────────────────────────────────

type DevTask struct {
	gorm.Model
	Title    string `gorm:"not null"`
	Type     string `gorm:"default:'Improvement'"`
	Assignee string
	Status   string `gorm:"default:'Todo'"`
	Priority string `gorm:"default:'Medium'"`
}

func (DevTask) TableName() string { return "dev_tasks" }
