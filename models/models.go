package models

import "gorm.io/gorm"

// ─── Sprint ───────────────────────────────────────────────────────────────────

type Sprint struct {
	gorm.Model
	ProjectID uint   `gorm:"index"`
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
	SprintID    uint                `gorm:"not null;index"`
	Title       string              `gorm:"not null"`
	AssigneeCSV string              `gorm:"column:assignees"`
	Status      string              `gorm:"default:'Todo'"`
	Priority    string              `gorm:"default:'Medium'"`
	Comments    []SprintTaskComment `gorm:"foreignKey:TaskID"`
	// Computed by service — not persisted
	Assignees []string `gorm:"-"`
}

type SprintTaskComment struct {
	gorm.Model
	TaskID  uint   `gorm:"not null;index"`
	Author  string `gorm:"not null"`
	Content string `gorm:"not null"`
}

func (SprintTaskComment) TableName() string { return "sprint_task_comments" }

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
	ProjectID uint   `gorm:"index"`
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
	ProjectID  uint   `gorm:"index"`
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
	ProjectID   uint   `gorm:"index"`
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
	ProjectID   uint             `gorm:"index"`
	Title       string           `gorm:"not null"`
	Type        string           `gorm:"default:'Improvement'"`
	AssigneeCSV string           `gorm:"column:assignees"`
	Status      string           `gorm:"default:'Todo'"`
	Priority    string           `gorm:"default:'Medium'"`
	Comments    []DevTaskComment `gorm:"foreignKey:TaskID"`
	// Computed by service — not persisted
	Assignees []string `gorm:"-"`
}

func (DevTask) TableName() string { return "dev_tasks" }

type DevTaskComment struct {
	gorm.Model
	TaskID  uint   `gorm:"not null;index"`
	Author  string `gorm:"not null"`
	Content string `gorm:"not null"`
}

func (DevTaskComment) TableName() string { return "dev_task_comments" }

// ─── Release ──────────────────────────────────────────────────────────────────

type Release struct {
	gorm.Model
	ProjectID   uint           `gorm:"index"`
	Name        string         `gorm:"not null"`
	Description string
	Status      string         `gorm:"default:'Draft'"` // Draft, In Progress, Released, Rolled Back
	TargetDate  string
	Stages      []ReleaseStage `gorm:"foreignKey:ReleaseID"`
}

type ReleaseStage struct {
	gorm.Model
	ReleaseID    uint                 `gorm:"not null;index"`
	Name         string               `gorm:"not null"`
	Status       string               `gorm:"default:'Pending'"` // Pending, Active, Done, Failed
	Stories      []ReleaseStory       `gorm:"foreignKey:StageID"`
	SlackUpdates []ReleaseSlackUpdate `gorm:"foreignKey:StageID"`
}

func (ReleaseStage) TableName() string { return "release_stages" }

type ReleaseStory struct {
	gorm.Model
	StageID  uint   `gorm:"not null;index"`
	Title    string `gorm:"not null"`
	Assignee string
	Status   string `gorm:"default:'Pending'"` // Pending, In QA, Passed, Failed
}

func (ReleaseStory) TableName() string { return "release_stories" }

type ReleaseSlackUpdate struct {
	gorm.Model
	StageID  uint   `gorm:"not null;index"`
	Channel  string
	Message  string `gorm:"not null"`
	Author   string
	PostedAt string
}

func (ReleaseSlackUpdate) TableName() string { return "release_slack_updates" }

// ─── Project & Team ───────────────────────────────────────────────────────────

type Project struct {
	gorm.Model
	Name        string       `gorm:"not null;uniqueIndex"`
	Description string
	Members     []TeamMember `gorm:"many2many:project_members;"`
}

type TeamMember struct {
	gorm.Model
	Name     string    `gorm:"not null"`
	Role     string
	Email    string
	Projects []Project `gorm:"many2many:project_members;"`
}

func (TeamMember) TableName() string { return "team_members" }
