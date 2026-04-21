package models

import "gorm.io/gorm"

// ─── User ─────────────────────────────────────────────────────────────────────

type User struct {
	gorm.Model
	FullName     string `gorm:"not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
}

func (User) TableName() string { return "users" }

// ─── Sprint (absorbs Release) ─────────────────────────────────────────────────

type Sprint struct {
	gorm.Model
	ProjectID   uint           `gorm:"index"`
	Name        string         `gorm:"not null"`
	Goal        string         `gorm:"default:''"`
	Progress    int            `gorm:"default:0"`
	StartDate   string
	EndDate     string
	Active      bool           `gorm:"default:false;index"`
	// Release fields
	Description string
	Status      string         `gorm:"default:'Draft'"` // Draft, In Progress, Released, Rolled Back
	TargetDate  string
	// Associations
	Tasks  []Task         `gorm:"foreignKey:SprintID"`
	Stages []ReleaseStage `gorm:"foreignKey:SprintID"`
}

type SprintStats struct {
	Todo       int
	InProgress int
	Done       int
	Blocked    int
}

func ComputeStats(tasks []Task) SprintStats {
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

// ─── Task (unified: sprint task / dev task / release task) ────────────────────

type Task struct {
	gorm.Model
	Category string `gorm:"not null;index"` // "sprint", "dev", "release"
	Title    string `gorm:"not null"`
	Status   string `gorm:"default:'Todo'"`
	Priority string `gorm:"default:'Medium'"`
	Type     string `gorm:"default:''"` // dev tasks: Improvement, Tech Debt, Research, Other

	// Parent FKs — only one is non-zero per task
	SprintID       uint `gorm:"index"`
	ProjectID      uint `gorm:"index"` // dev tasks
	ReleaseStageID uint `gorm:"index"` // release tasks

	// Release tasks use a single assignee; sprint/dev use Assignees (many2many)
	AssigneeID *uint       `gorm:"index"`
	Assignee   *TeamMember `gorm:"foreignKey:AssigneeID"`

	Comments  []TaskComment `gorm:"foreignKey:TaskID"`
	Assignees []TeamMember  `gorm:"many2many:task_assignees;"`
}

func (Task) TableName() string { return "tasks" }

type TaskComment struct {
	gorm.Model
	TaskID  uint   `gorm:"not null;index"`
	Author  string `gorm:"not null"`
	Content string `gorm:"not null"`
}

func (TaskComment) TableName() string { return "task_comments" }

// ─── Release Stage ────────────────────────────────────────────────────────────

type ReleaseStage struct {
	gorm.Model
	SprintID     uint                 `gorm:"index"` // FK → Sprint.ID
	Name         string               `gorm:"not null"`
	Status       string               `gorm:"default:'Pending'"` // Pending, Active, Done, Failed
	Stories      []Task               `gorm:"foreignKey:ReleaseStageID"`
	SlackUpdates []ReleaseSlackUpdate `gorm:"foreignKey:StageID"`
}

func (ReleaseStage) TableName() string { return "release_stages" }

type ReleaseSlackUpdate struct {
	gorm.Model
	StageID  uint   `gorm:"not null;index"`
	Channel  string
	Message  string `gorm:"not null"`
	Author   string
	PostedAt string
}

func (ReleaseSlackUpdate) TableName() string { return "release_slack_updates" }

// ─── Standup ──────────────────────────────────────────────────────────────────

type StandupEntry struct {
	gorm.Model
	ProjectID    uint    `gorm:"index"`
	Project      Project `gorm:"foreignKey:ProjectID"`
	Date         string  `gorm:"not null;index"` // "YYYY-MM-DD"
	Summary      string  `gorm:"type:text"`
	Dependencies string  `gorm:"type:text"`
	Issues       string  `gorm:"type:text"`
	ActionItems  string  `gorm:"type:text"`
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
	ProjectID   uint    `gorm:"index"`
	Project     Project `gorm:"foreignKey:ProjectID"`
	Title       string  `gorm:"not null"`
	Date        string  `gorm:"not null"`
	AttendeeCSV string  `gorm:"column:attendees"`
	Notes       string
	ActionItems []ActionItem `gorm:"foreignKey:MeetingID"`
	// Computed by service — not persisted
	Attendees []string `gorm:"-"`
}

// ─── Sticky Note ──────────────────────────────────────────────────────────────

type StickyNote struct {
	gorm.Model
	ProjectID uint   `gorm:"index"`
	Title     string
	Content   string `gorm:"type:text"`
	Color     string `gorm:"default:'yellow'"` // yellow, green, blue, pink, purple
	Pinned    bool   `gorm:"default:false"`
}

func (StickyNote) TableName() string { return "sticky_notes" }

// ─── Slack Thread ─────────────────────────────────────────────────────────────

type SlackThread struct {
	gorm.Model
	ProjectID   uint        `gorm:"index"`
	MessageLink string
	Topic       string      `gorm:"not null"`
	Summary     string
	TagCSV      string      `gorm:"column:tags"`
	AuthorID    *uint       `gorm:"index"`
	Author      *TeamMember `gorm:"foreignKey:AuthorID"`
	// Computed by service — not persisted
	Tags []string `gorm:"-"`
}

func (SlackThread) TableName() string { return "slack_threads" }

// ─── Project & Team ───────────────────────────────────────────────────────────

type Project struct {
	gorm.Model
	Name        string       `gorm:"not null;uniqueIndex"`
	Description string
	Members     []TeamMember `gorm:"many2many:project_members;"`
}

type TeamMember struct {
	gorm.Model
	UserID   *uint     `gorm:"index"`
	User     *User     `gorm:"foreignKey:UserID"`
	Name     string    `gorm:"not null"`
	Role     string
	Email    string
	Projects []Project `gorm:"many2many:project_members;"`
}

func (TeamMember) TableName() string { return "team_members" }
