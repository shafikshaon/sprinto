package db

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"sprinto/models"
)

func Connect(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("db.Connect: %v\nMake sure PostgreSQL is running.\nCreate the DB with: createdb sprinto", err)
	}
	return db
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Sprint{},
		&models.SprintTask{},
		&models.StandupEntry{},
		&models.Deadline{},
		&models.Meeting{},
		&models.ActionItem{},
		&models.DevTask{},
	)
}

func IsEmpty(db *gorm.DB) bool {
	var count int64
	db.Model(&models.Sprint{}).Count(&count)
	return count == 0
}

func Seed(db *gorm.DB) error {
	sprint := &models.Sprint{
		Name:      "Sprint 12",
		Goal:      "Ship OAuth2 integration and complete API v2 endpoints",
		Progress:  62,
		StartDate: "Apr 14",
		EndDate:   "Apr 28",
		Active:    true,
	}
	if err := db.Create(sprint).Error; err != nil {
		return err
	}

	tasks := []models.SprintTask{
		{SprintID: sprint.ID, Title: "Fix authentication token refresh", Assignee: "Alice Chen", Status: "In Progress", Priority: "High"},
		{SprintID: sprint.ID, Title: "Implement OAuth2 provider", Assignee: "Bob Martinez", Status: "In Progress", Priority: "High"},
		{SprintID: sprint.ID, Title: "Database migration script", Assignee: "Carol Singh", Status: "Done", Priority: "Medium"},
		{SprintID: sprint.ID, Title: "CI/CD pipeline optimization", Assignee: "Dan Kim", Status: "Todo", Priority: "Medium"},
		{SprintID: sprint.ID, Title: "Write E2E test suite", Assignee: "Eva Park", Status: "Todo", Priority: "High"},
		{SprintID: sprint.ID, Title: "API rate limiting", Assignee: "Alice Chen", Status: "Done", Priority: "Medium"},
		{SprintID: sprint.ID, Title: "Mobile responsive fixes", Assignee: "Bob Martinez", Status: "Todo", Priority: "Low"},
		{SprintID: sprint.ID, Title: "Cache layer implementation", Assignee: "Carol Singh", Status: "In Progress", Priority: "High"},
		{SprintID: sprint.ID, Title: "Docker Compose setup", Assignee: "Dan Kim", Status: "Done", Priority: "Low"},
		{SprintID: sprint.ID, Title: "Performance regression tests", Assignee: "Eva Park", Status: "In Progress", Priority: "Medium"},
	}
	if err := db.Create(&tasks).Error; err != nil {
		return err
	}

	today := time.Now().Format("2006-01-02")
	standups := []models.StandupEntry{
		{Member: "Alice Chen", Role: "Backend", Yesterday: "Fixed JWT expiry bug in auth middleware", Today: "Continue OAuth2 integration with provider SDK", Blockers: "None", Status: "On Track", Date: today},
		{Member: "Bob Martinez", Role: "Frontend", Yesterday: "Updated UI component library to v3", Today: "Implement dashboard charts and sprint progress view", Blockers: "Waiting on design approval for new layout", Status: "At Risk", Date: today},
		{Member: "Carol Singh", Role: "Full Stack", Yesterday: "Completed DB migration script and tested in staging", Today: "Start cache layer implementation with Redis", Blockers: "None", Status: "On Track", Date: today},
		{Member: "Dan Kim", Role: "DevOps", Yesterday: "Reviewed Kubernetes config and updated resource limits", Today: "Optimize CI/CD pipeline run times", Blockers: "Need access to production cluster", Status: "Blocked", Date: today},
		{Member: "Eva Park", Role: "QA", Yesterday: "Created test plan for sprint 12 features", Today: "Set up Playwright E2E test framework", Blockers: "None", Status: "On Track", Date: today},
	}
	if err := db.Create(&standups).Error; err != nil {
		return err
	}

	deadlines := []models.Deadline{
		{Title: "API v2 Launch", Project: "Platform", DueDateRaw: "2026-04-20", Priority: "Critical"},
		{Title: "Mobile App Beta", Project: "Mobile", DueDateRaw: "2026-04-25", Priority: "High"},
		{Title: "Security Audit", Project: "Security", DueDateRaw: "2026-04-30", Priority: "High"},
		{Title: "Q2 OKR Review", Project: "Management", DueDateRaw: "2026-05-05", Priority: "Medium"},
		{Title: "Infrastructure Upgrade", Project: "Infrastructure", DueDateRaw: "2026-05-15", Priority: "Low"},
	}
	if err := db.Create(&deadlines).Error; err != nil {
		return err
	}

	meetings := []models.Meeting{
		{Title: "Sprint 12 Planning", Date: "Apr 14, 2026", AttendeeCSV: "Alice Chen,Bob Martinez,Carol Singh,Dan Kim,Eva Park", Notes: "Defined sprint goals focused on OAuth2 and API v2. Capacity at 85% due to Dan's infrastructure work. Story points committed: 42."},
		{Title: "1:1 — Alice Chen", Date: "Apr 15, 2026", AttendeeCSV: "Alice Chen", Notes: "Discussed career growth path toward Staff Engineer. Alice interested in leading API v2 architecture initiative next quarter."},
		{Title: "Architecture Review — API v2", Date: "Apr 16, 2026", AttendeeCSV: "Alice Chen,Bob Martinez,Carol Singh", Notes: "Decided to keep REST with improved versioning strategy and OpenAPI docs."},
		{Title: "Incident Post-mortem", Date: "Apr 17, 2026", AttendeeCSV: "Alice Chen,Bob Martinez,Carol Singh,Dan Kim,Eva Park", Notes: "Auth service outage Apr 16 (45 min). Root cause: TLS cert expired without alert. No data loss. SLA impacted."},
	}
	if err := db.Create(&meetings).Error; err != nil {
		return err
	}
	db.Create(&models.ActionItem{MeetingID: meetings[0].ID, Task: "Share OAuth2 spec doc with team", Owner: "Alice Chen", DueDate: "Apr 15", Done: true})
	db.Create(&models.ActionItem{MeetingID: meetings[0].ID, Task: "Update design mockups for dashboard", Owner: "Bob Martinez", DueDate: "Apr 16", Done: false})
	db.Create(&models.ActionItem{MeetingID: meetings[3].ID, Task: "Set up cert expiry monitoring alerts", Owner: "Dan Kim", DueDate: "Apr 20", Done: false})
	db.Create(&models.ActionItem{MeetingID: meetings[3].ID, Task: "Write incident summary for stakeholders", Owner: "EM", DueDate: "Apr 18", Done: false})

	devTasks := []models.DevTask{
		{Title: "Migrate REST endpoints to OpenAPI spec", Type: "Improvement", Assignee: "Alice Chen", Status: "In Progress", Priority: "High"},
		{Title: "Upgrade Go version to 1.22", Type: "Tech Debt", Assignee: "Dan Kim", Status: "In Progress", Priority: "Medium"},
		{Title: "Add Swagger UI for API documentation", Type: "Improvement", Assignee: "Bob Martinez", Status: "Todo", Priority: "Medium"},
		{Title: "Remove deprecated auth middleware", Type: "Tech Debt", Assignee: "Carol Singh", Status: "Done", Priority: "High"},
		{Title: "Benchmark critical database queries", Type: "Research", Assignee: "Alice Chen", Status: "Todo", Priority: "Low"},
		{Title: "Setup Sentry error monitoring", Type: "Improvement", Assignee: "Dan Kim", Status: "Todo", Priority: "High"},
		{Title: "Consolidate app config management", Type: "Tech Debt", Assignee: "Carol Singh", Status: "Todo", Priority: "Medium"},
		{Title: "Load testing with k6", Type: "Research", Assignee: "Eva Park", Status: "Todo", Priority: "Medium"},
	}
	return db.Create(&devTasks).Error
}
