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
	if err := db.AutoMigrate(
		&models.User{},
		&models.Sprint{},
		&models.SprintTask{},
		&models.SprintTaskComment{},
		&models.StandupEntry{},
		&models.Deadline{},
		&models.Meeting{},
		&models.ActionItem{},
		&models.DevTask{},
		&models.DevTaskComment{},
		&models.Release{},
		&models.ReleaseStage{},
		&models.ReleaseStory{},
		&models.ReleaseSlackUpdate{},
		&models.Project{},
		&models.TeamMember{},
	); err != nil {
		return err
	}
	// Copy old single-assignee column into new multi-assignee CSV column (idempotent)
	db.Exec("UPDATE sprint_tasks SET assignees = assignee WHERE (assignees IS NULL OR assignees = '') AND assignee IS NOT NULL AND assignee != ''")
	db.Exec("UPDATE dev_tasks SET assignees = assignee WHERE (assignees IS NULL OR assignees = '') AND assignee IS NOT NULL AND assignee != ''")
	return nil
}

func IsEmpty(db *gorm.DB) bool {
	var count int64
	db.Model(&models.Project{}).Count(&count)
	return count == 0
}

func Seed(db *gorm.DB) error {
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	twoDaysAgo := time.Now().AddDate(0, 0, -2).Format("2006-01-02")

	// ── Team Members ─────────────────────────────────────────────────────────
	alice := &models.TeamMember{Name: "Alice Chen", Role: "Backend Engineer", Email: "alice@example.com"}
	bob := &models.TeamMember{Name: "Bob Martinez", Role: "Frontend Engineer", Email: "bob@example.com"}
	carol := &models.TeamMember{Name: "Carol Singh", Role: "Full Stack Engineer", Email: "carol@example.com"}
	dan := &models.TeamMember{Name: "Dan Kim", Role: "DevOps Engineer", Email: "dan@example.com"}
	eva := &models.TeamMember{Name: "Eva Park", Role: "QA Engineer", Email: "eva@example.com"}
	frank := &models.TeamMember{Name: "Frank Liu", Role: "iOS Engineer", Email: "frank@example.com"}
	grace := &models.TeamMember{Name: "Grace Obi", Role: "Android Engineer", Email: "grace@example.com"}
	henry := &models.TeamMember{Name: "Henry Walsh", Role: "Security Engineer", Email: "henry@example.com"}
	iris := &models.TeamMember{Name: "Iris Novak", Role: "Compliance Analyst", Email: "iris@example.com"}
	for _, m := range []*models.TeamMember{alice, bob, carol, dan, eva, frank, grace, henry, iris} {
		if err := db.Create(m).Error; err != nil {
			return err
		}
	}

	// ── Projects ─────────────────────────────────────────────────────────────
	platform := &models.Project{Name: "Platform", Description: "Core API and backend infrastructure services"}
	mobile := &models.Project{Name: "Mobile", Description: "iOS and Android applications"}
	security := &models.Project{Name: "Security", Description: "Auth, compliance, and security hardening"}
	for _, p := range []*models.Project{platform, mobile, security} {
		if err := db.Create(p).Error; err != nil {
			return err
		}
	}
	db.Model(platform).Association("Members").Append(alice, bob, carol, dan, eva)
	db.Model(mobile).Association("Members").Append(bob, carol, frank, grace, eva)
	db.Model(security).Association("Members").Append(alice, dan, henry, iris)

	// ═══════════════════════════════════════════════════════════════════════════
	// PLATFORM PROJECT
	// ═══════════════════════════════════════════════════════════════════════════

	// ── Platform Sprint ───────────────────────────────────────────────────────
	platformSprint := &models.Sprint{
		ProjectID: platform.ID,
		Name:      "Sprint 12",
		Goal:      "Ship OAuth2 integration and complete API v2 endpoints",
		Progress:  62,
		StartDate: "Apr 14",
		EndDate:   "Apr 28",
		Active:    true,
	}
	if err := db.Create(platformSprint).Error; err != nil {
		return err
	}
	db.Create(&[]models.SprintTask{
		{SprintID: platformSprint.ID, Title: "Fix authentication token refresh", AssigneeCSV: "Alice Chen", Status: "In Progress", Priority: "High"},
		{SprintID: platformSprint.ID, Title: "Implement OAuth2 provider", AssigneeCSV: "Bob Martinez", Status: "In Progress", Priority: "High"},
		{SprintID: platformSprint.ID, Title: "Database migration script", AssigneeCSV: "Carol Singh", Status: "Done", Priority: "Medium"},
		{SprintID: platformSprint.ID, Title: "CI/CD pipeline optimisation", AssigneeCSV: "Dan Kim", Status: "Todo", Priority: "Medium"},
		{SprintID: platformSprint.ID, Title: "Write E2E test suite", AssigneeCSV: "Eva Park", Status: "Todo", Priority: "High"},
		{SprintID: platformSprint.ID, Title: "API rate limiting", AssigneeCSV: "Alice Chen", Status: "Done", Priority: "Medium"},
		{SprintID: platformSprint.ID, Title: "Mobile responsive fixes", AssigneeCSV: "Bob Martinez", Status: "Todo", Priority: "Low"},
		{SprintID: platformSprint.ID, Title: "Cache layer implementation", AssigneeCSV: "Carol Singh", Status: "In Progress", Priority: "High"},
		{SprintID: platformSprint.ID, Title: "Docker Compose setup", AssigneeCSV: "Dan Kim", Status: "Done", Priority: "Low"},
		{SprintID: platformSprint.ID, Title: "Performance regression tests", AssigneeCSV: "Eva Park", Status: "In Progress", Priority: "Medium"},
	})

	// ── Platform Standups ─────────────────────────────────────────────────────
	db.Create(&[]models.StandupEntry{
		{ProjectID: platform.ID, Member: "Alice Chen", Role: "Backend", Yesterday: "Fixed JWT expiry bug in auth middleware", Today: "Continue OAuth2 integration with provider SDK", Blockers: "None", Status: "On Track", Date: today},
		{ProjectID: platform.ID, Member: "Bob Martinez", Role: "Frontend", Yesterday: "Updated UI component library to v3", Today: "Implement dashboard charts and sprint progress view", Blockers: "Waiting on design approval for new layout", Status: "At Risk", Date: today},
		{ProjectID: platform.ID, Member: "Carol Singh", Role: "Full Stack", Yesterday: "Completed DB migration script and tested in staging", Today: "Start cache layer implementation with Redis", Blockers: "None", Status: "On Track", Date: today},
		{ProjectID: platform.ID, Member: "Dan Kim", Role: "DevOps", Yesterday: "Reviewed Kubernetes config and updated resource limits", Today: "Optimise CI/CD pipeline run times", Blockers: "Need access to production cluster", Status: "Blocked", Date: today},
		{ProjectID: platform.ID, Member: "Eva Park", Role: "QA", Yesterday: "Created test plan for Sprint 12 features", Today: "Set up Playwright E2E test framework", Blockers: "None", Status: "On Track", Date: today},
		// yesterday
		{ProjectID: platform.ID, Member: "Alice Chen", Role: "Backend", Yesterday: "Reviewed OAuth2 spec and opened provider account", Today: "Start token refresh fix", Blockers: "None", Status: "On Track", Date: yesterday},
		{ProjectID: platform.ID, Member: "Bob Martinez", Role: "Frontend", Yesterday: "Audited existing UI components", Today: "Upgrade component library", Blockers: "None", Status: "On Track", Date: yesterday},
		{ProjectID: platform.ID, Member: "Carol Singh", Role: "Full Stack", Yesterday: "Wrote DB migration plan", Today: "Implement and test migration script", Blockers: "None", Status: "On Track", Date: yesterday},
		// two days ago
		{ProjectID: platform.ID, Member: "Alice Chen", Role: "Backend", Yesterday: "Sprint planning and task breakdown", Today: "Begin OAuth2 research", Blockers: "None", Status: "On Track", Date: twoDaysAgo},
		{ProjectID: platform.ID, Member: "Dan Kim", Role: "DevOps", Yesterday: "Upgraded staging K8s cluster", Today: "Review prod cluster config", Blockers: "None", Status: "On Track", Date: twoDaysAgo},
	})

	// ── Platform Deadlines ────────────────────────────────────────────────────
	db.Create(&[]models.Deadline{
		{ProjectID: platform.ID, Title: "API v2 Public Launch", Project: "Platform", DueDateRaw: "2026-04-20", Priority: "Critical"},
		{ProjectID: platform.ID, Title: "OAuth2 Integration Complete", Project: "Platform", DueDateRaw: "2026-04-28", Priority: "High"},
		{ProjectID: platform.ID, Title: "Q2 Engineering OKR Review", Project: "Platform", DueDateRaw: "2026-05-05", Priority: "Medium"},
		{ProjectID: platform.ID, Title: "Infrastructure Upgrade — Phase 2", Project: "Platform", DueDateRaw: "2026-05-20", Priority: "Low"},
	})

	// ── Platform Meetings ─────────────────────────────────────────────────────
	pm1 := &models.Meeting{ProjectID: platform.ID, Title: "Sprint 12 Planning", Date: "Apr 14, 2026", AttendeeCSV: "Alice Chen,Bob Martinez,Carol Singh,Dan Kim,Eva Park", Notes: "Defined sprint goals focused on OAuth2 and API v2. Capacity at 85% due to Dan's infrastructure work. Story points committed: 42."}
	pm2 := &models.Meeting{ProjectID: platform.ID, Title: "1:1 — Alice Chen", Date: "Apr 15, 2026", AttendeeCSV: "Alice Chen", Notes: "Discussed career growth path toward Staff Engineer. Alice interested in leading API v2 architecture initiative next quarter."}
	pm3 := &models.Meeting{ProjectID: platform.ID, Title: "Architecture Review — API v2", Date: "Apr 16, 2026", AttendeeCSV: "Alice Chen,Bob Martinez,Carol Singh", Notes: "Decided to keep REST with improved versioning strategy and OpenAPI docs. GraphQL deferred to Q3."}
	pm4 := &models.Meeting{ProjectID: platform.ID, Title: "Incident Post-mortem", Date: "Apr 17, 2026", AttendeeCSV: "Alice Chen,Bob Martinez,Carol Singh,Dan Kim,Eva Park", Notes: "Auth service outage Apr 16 (45 min). Root cause: TLS cert expired without alert. No data loss. SLA impacted."}
	for _, m := range []*models.Meeting{pm1, pm2, pm3, pm4} {
		db.Create(m)
	}
	db.Create(&models.ActionItem{MeetingID: pm1.ID, Task: "Share OAuth2 spec doc with team", Owner: "Alice Chen", DueDate: "Apr 15", Done: true})
	db.Create(&models.ActionItem{MeetingID: pm1.ID, Task: "Update design mockups for dashboard", Owner: "Bob Martinez", DueDate: "Apr 16", Done: false})
	db.Create(&models.ActionItem{MeetingID: pm3.ID, Task: "Write OpenAPI schema for v2 endpoints", Owner: "Alice Chen", DueDate: "Apr 21", Done: false})
	db.Create(&models.ActionItem{MeetingID: pm4.ID, Task: "Set up cert expiry monitoring alerts", Owner: "Dan Kim", DueDate: "Apr 20", Done: false})
	db.Create(&models.ActionItem{MeetingID: pm4.ID, Task: "Write incident summary for stakeholders", Owner: "EM", DueDate: "Apr 18", Done: false})

	// ── Platform Dev Tasks ────────────────────────────────────────────────────
	db.Create(&[]models.DevTask{
		{ProjectID: platform.ID, Title: "Migrate REST endpoints to OpenAPI spec", Type: "Improvement", AssigneeCSV: "Alice Chen", Status: "In Progress", Priority: "High"},
		{ProjectID: platform.ID, Title: "Upgrade Go version to 1.22", Type: "Tech Debt", AssigneeCSV: "Dan Kim", Status: "In Progress", Priority: "Medium"},
		{ProjectID: platform.ID, Title: "Add Swagger UI for API documentation", Type: "Improvement", AssigneeCSV: "Bob Martinez", Status: "Todo", Priority: "Medium"},
		{ProjectID: platform.ID, Title: "Remove deprecated auth middleware", Type: "Tech Debt", AssigneeCSV: "Carol Singh", Status: "Done", Priority: "High"},
		{ProjectID: platform.ID, Title: "Benchmark critical database queries", Type: "Research", AssigneeCSV: "Alice Chen", Status: "Todo", Priority: "Low"},
		{ProjectID: platform.ID, Title: "Setup Sentry error monitoring", Type: "Improvement", AssigneeCSV: "Dan Kim", Status: "Todo", Priority: "High"},
		{ProjectID: platform.ID, Title: "Consolidate app config management", Type: "Tech Debt", AssigneeCSV: "Carol Singh", Status: "Todo", Priority: "Medium"},
		{ProjectID: platform.ID, Title: "Load testing with k6", Type: "Research", AssigneeCSV: "Eva Park", Status: "Todo", Priority: "Medium"},
	})

	// ── Platform Releases ─────────────────────────────────────────────────────
	rel1 := &models.Release{ProjectID: platform.ID, Name: "v2.3.0 – Auth Overhaul", Description: "OAuth2 integration and JWT improvements", Status: "In Progress", TargetDate: "2026-04-28"}
	db.Create(rel1)
	s1 := &models.ReleaseStage{ReleaseID: rel1.ID, Name: "QA Round 1", Status: "Done"}
	s2 := &models.ReleaseStage{ReleaseID: rel1.ID, Name: "QA Round 2", Status: "Active"}
	s3 := &models.ReleaseStage{ReleaseID: rel1.ID, Name: "Staging Deploy", Status: "Pending"}
	db.Create(s1); db.Create(s2); db.Create(s3)
	db.Create(&models.ReleaseStory{StageID: s1.ID, Title: "Fix authentication token refresh", Assignee: "Eva Park", Status: "Passed"})
	db.Create(&models.ReleaseStory{StageID: s1.ID, Title: "API rate limiting", Assignee: "Eva Park", Status: "Passed"})
	db.Create(&models.ReleaseSlackUpdate{StageID: s1.ID, Channel: "#releases", Message: "QA Round 1 complete — all stories passed. Promoting to Round 2.", Author: "Eva Park", PostedAt: "Apr 16, 2:30 PM"})
	db.Create(&models.ReleaseStory{StageID: s2.ID, Title: "Implement OAuth2 provider", Assignee: "Eva Park", Status: "In QA"})
	db.Create(&models.ReleaseStory{StageID: s2.ID, Title: "Cache layer implementation", Assignee: "Eva Park", Status: "Pending"})
	db.Create(&models.ReleaseSlackUpdate{StageID: s2.ID, Channel: "#releases", Message: "Starting QA Round 2. OAuth2 story handed to Eva. Cache layer follows once OAuth2 passes.", Author: "Alice Chen", PostedAt: "Apr 17, 10:15 AM"})

	rel2 := &models.Release{ProjectID: platform.ID, Name: "v2.2.1 – Security Hotfix", Description: "TLS certificate rotation and auth session hardening", Status: "Released", TargetDate: "2026-04-18"}
	db.Create(rel2)
	hs1 := &models.ReleaseStage{ReleaseID: rel2.ID, Name: "QA Verification", Status: "Done"}
	hs2 := &models.ReleaseStage{ReleaseID: rel2.ID, Name: "Production Deploy", Status: "Done"}
	db.Create(hs1); db.Create(hs2)
	db.Create(&models.ReleaseStory{StageID: hs1.ID, Title: "TLS cert auto-rotation script", Assignee: "Dan Kim", Status: "Passed"})
	db.Create(&models.ReleaseStory{StageID: hs1.ID, Title: "Session token TTL enforcement", Assignee: "Alice Chen", Status: "Passed"})
	db.Create(&models.ReleaseSlackUpdate{StageID: hs1.ID, Channel: "#releases", Message: "Hotfix verified in staging. Both stories passed QA. Requesting prod deploy window.", Author: "Eva Park", PostedAt: "Apr 17, 4:00 PM"})
	db.Create(&models.ReleaseSlackUpdate{StageID: hs2.ID, Channel: "#incidents", Message: "v2.2.1 deployed to prod. TLS rotation confirmed. Monitoring for 30 min.", Author: "Dan Kim", PostedAt: "Apr 17, 6:45 PM"})
	db.Create(&models.ReleaseSlackUpdate{StageID: hs2.ID, Channel: "#releases", Message: "v2.2.1 stable — no issues. Release complete.", Author: "EM", PostedAt: "Apr 17, 7:20 PM"})

	// ═══════════════════════════════════════════════════════════════════════════
	// MOBILE PROJECT
	// ═══════════════════════════════════════════════════════════════════════════

	mobileSprint := &models.Sprint{
		ProjectID: mobile.ID,
		Name:      "Sprint 8",
		Goal:      "Launch push notifications and implement offline mode for core screens",
		Progress:  41,
		StartDate: "Apr 14",
		EndDate:   "Apr 28",
		Active:    true,
	}
	if err := db.Create(mobileSprint).Error; err != nil {
		return err
	}
	db.Create(&[]models.SprintTask{
		{SprintID: mobileSprint.ID, Title: "Push notification service integration", AssigneeCSV: "Frank Liu", Status: "In Progress", Priority: "High"},
		{SprintID: mobileSprint.ID, Title: "Offline data sync for home screen", AssigneeCSV: "Grace Obi", Status: "In Progress", Priority: "High"},
		{SprintID: mobileSprint.ID, Title: "iOS deep link handling", AssigneeCSV: "Frank Liu", Status: "Todo", Priority: "Medium"},
		{SprintID: mobileSprint.ID, Title: "Android background fetch", AssigneeCSV: "Grace Obi", Status: "Todo", Priority: "Medium"},
		{SprintID: mobileSprint.ID, Title: "Notification permission flow UI", AssigneeCSV: "Bob Martinez", Status: "Done", Priority: "High"},
		{SprintID: mobileSprint.ID, Title: "SQLite offline schema migration", AssigneeCSV: "Carol Singh", Status: "In Progress", Priority: "High"},
		{SprintID: mobileSprint.ID, Title: "E2E tests for notification scenarios", AssigneeCSV: "Eva Park", Status: "Todo", Priority: "Medium"},
		{SprintID: mobileSprint.ID, Title: "App Store release notes draft", AssigneeCSV: "Bob Martinez", Status: "Todo", Priority: "Low"},
	})

	// ── Mobile Standups ───────────────────────────────────────────────────────
	db.Create(&[]models.StandupEntry{
		{ProjectID: mobile.ID, Member: "Frank Liu", Role: "iOS", Yesterday: "Integrated APNs token registration", Today: "Handle notification payload parsing and display", Blockers: "None", Status: "On Track", Date: today},
		{ProjectID: mobile.ID, Member: "Grace Obi", Role: "Android", Yesterday: "Set up WorkManager for background sync", Today: "Implement conflict resolution for offline data sync", Blockers: "FCM quota limit hit in dev — waiting for increase", Status: "At Risk", Date: today},
		{ProjectID: mobile.ID, Member: "Bob Martinez", Role: "Frontend", Yesterday: "Shipped notification permission UI screens", Today: "Polish onboarding flow animations", Blockers: "None", Status: "On Track", Date: today},
		{ProjectID: mobile.ID, Member: "Carol Singh", Role: "Full Stack", Yesterday: "Wrote SQLite schema for offline tables", Today: "Run migration on device simulators and fix edge cases", Blockers: "None", Status: "On Track", Date: today},
		{ProjectID: mobile.ID, Member: "Eva Park", Role: "QA", Yesterday: "Explored notification test tooling", Today: "Write test cases for notification permission scenarios", Blockers: "None", Status: "On Track", Date: today},
		// yesterday
		{ProjectID: mobile.ID, Member: "Frank Liu", Role: "iOS", Yesterday: "Set up APNs certificates", Today: "Integrate APNs token registration", Blockers: "None", Status: "On Track", Date: yesterday},
		{ProjectID: mobile.ID, Member: "Grace Obi", Role: "Android", Yesterday: "Researched WorkManager API", Today: "Set up WorkManager for background sync", Blockers: "None", Status: "On Track", Date: yesterday},
	})

	// ── Mobile Deadlines ──────────────────────────────────────────────────────
	db.Create(&[]models.Deadline{
		{ProjectID: mobile.ID, Title: "App Store Submission — v1.2", Project: "Mobile", DueDateRaw: "2026-04-25", Priority: "Critical"},
		{ProjectID: mobile.ID, Title: "Beta TestFlight Release", Project: "Mobile", DueDateRaw: "2026-04-22", Priority: "High"},
		{ProjectID: mobile.ID, Title: "QA Sign-off for Push Notifications", Project: "Mobile", DueDateRaw: "2026-04-21", Priority: "High"},
		{ProjectID: mobile.ID, Title: "Google Play Internal Testing", Project: "Mobile", DueDateRaw: "2026-05-01", Priority: "Medium"},
	})

	// ── Mobile Meetings ───────────────────────────────────────────────────────
	mm1 := &models.Meeting{ProjectID: mobile.ID, Title: "Sprint 8 Kickoff", Date: "Apr 14, 2026", AttendeeCSV: "Frank Liu,Grace Obi,Bob Martinez,Carol Singh,Eva Park", Notes: "Aligned on push notification architecture. Decided to use Firebase for Android and APNs for iOS with a shared backend abstraction layer."}
	mm2 := &models.Meeting{ProjectID: mobile.ID, Title: "App Store Review Prep", Date: "Apr 16, 2026", AttendeeCSV: "Frank Liu,Bob Martinez", Notes: "Reviewed App Store guidelines for notification features. Privacy manifest needs updating — required for submission."}
	for _, m := range []*models.Meeting{mm1, mm2} {
		db.Create(m)
	}
	db.Create(&models.ActionItem{MeetingID: mm1.ID, Task: "Set up shared push notification abstraction in backend", Owner: "Carol Singh", DueDate: "Apr 17", Done: true})
	db.Create(&models.ActionItem{MeetingID: mm1.ID, Task: "Create Firebase project and share credentials", Owner: "Grace Obi", DueDate: "Apr 15", Done: true})
	db.Create(&models.ActionItem{MeetingID: mm2.ID, Task: "Update PrivacyInfo.xcprivacy manifest", Owner: "Frank Liu", DueDate: "Apr 19", Done: false})
	db.Create(&models.ActionItem{MeetingID: mm2.ID, Task: "Screenshot all new notification UI screens", Owner: "Bob Martinez", DueDate: "Apr 20", Done: false})

	// ── Mobile Dev Tasks ──────────────────────────────────────────────────────
	db.Create(&[]models.DevTask{
		{ProjectID: mobile.ID, Title: "Upgrade React Native to 0.74", Type: "Tech Debt", AssigneeCSV: "Bob Martinez", Status: "Todo", Priority: "High"},
		{ProjectID: mobile.ID, Title: "Implement biometric auth (Face ID / fingerprint)", Type: "Improvement", AssigneeCSV: "Frank Liu", Status: "Todo", Priority: "Medium"},
		{ProjectID: mobile.ID, Title: "Benchmark SQLite vs Realm for offline storage", Type: "Research", AssigneeCSV: "Carol Singh", Status: "In Progress", Priority: "Medium"},
		{ProjectID: mobile.ID, Title: "Remove legacy Bluetooth module", Type: "Tech Debt", AssigneeCSV: "Grace Obi", Status: "Done", Priority: "Low"},
		{ProjectID: mobile.ID, Title: "Add crash analytics with Crashlytics", Type: "Improvement", AssigneeCSV: "Frank Liu", Status: "Todo", Priority: "High"},
	})

	// ── Mobile Releases ───────────────────────────────────────────────────────
	mrel1 := &models.Release{ProjectID: mobile.ID, Name: "v1.2.0 – Push & Offline", Description: "Push notifications and offline mode for core screens", Status: "In Progress", TargetDate: "2026-04-25"}
	db.Create(mrel1)
	ms1 := &models.ReleaseStage{ReleaseID: mrel1.ID, Name: "Internal QA", Status: "Active"}
	ms2 := &models.ReleaseStage{ReleaseID: mrel1.ID, Name: "Beta — TestFlight / Play", Status: "Pending"}
	ms3 := &models.ReleaseStage{ReleaseID: mrel1.ID, Name: "App Store / Play Store", Status: "Pending"}
	db.Create(ms1); db.Create(ms2); db.Create(ms3)
	db.Create(&models.ReleaseStory{StageID: ms1.ID, Title: "Push notification permission flow", Assignee: "Eva Park", Status: "Passed"})
	db.Create(&models.ReleaseStory{StageID: ms1.ID, Title: "Push notification delivery — iOS", Assignee: "Eva Park", Status: "In QA"})
	db.Create(&models.ReleaseStory{StageID: ms1.ID, Title: "Push notification delivery — Android", Assignee: "Eva Park", Status: "Pending"})
	db.Create(&models.ReleaseSlackUpdate{StageID: ms1.ID, Channel: "#mobile-releases", Message: "Internal QA started. Permission flow passed. iOS push delivery testing in progress.", Author: "Eva Park", PostedAt: "Apr 17, 9:00 AM"})

	// ═══════════════════════════════════════════════════════════════════════════
	// SECURITY PROJECT
	// ═══════════════════════════════════════════════════════════════════════════

	securitySprint := &models.Sprint{
		ProjectID: security.ID,
		Name:      "Sprint 5",
		Goal:      "Complete penetration test remediation and achieve SOC 2 Type II readiness",
		Progress:  28,
		StartDate: "Apr 7",
		EndDate:   "Apr 25",
		Active:    true,
	}
	if err := db.Create(securitySprint).Error; err != nil {
		return err
	}
	db.Create(&[]models.SprintTask{
		{SprintID: securitySprint.ID, Title: "Remediate pen test finding: SQL injection in search", AssigneeCSV: "Henry Walsh", Status: "Done", Priority: "Critical"},
		{SprintID: securitySprint.ID, Title: "Remediate pen test finding: IDOR in user API", AssigneeCSV: "Alice Chen", Status: "In Progress", Priority: "Critical"},
		{SprintID: securitySprint.ID, Title: "Implement audit log for all admin actions", AssigneeCSV: "Henry Walsh", Status: "In Progress", Priority: "High"},
		{SprintID: securitySprint.ID, Title: "Enable MFA enforcement for admin accounts", AssigneeCSV: "Dan Kim", Status: "Done", Priority: "High"},
		{SprintID: securitySprint.ID, Title: "Review and rotate all production secrets", AssigneeCSV: "Dan Kim", Status: "Todo", Priority: "High"},
		{SprintID: securitySprint.ID, Title: "Document access control matrix for SOC 2", AssigneeCSV: "Iris Novak", Status: "In Progress", Priority: "Medium"},
		{SprintID: securitySprint.ID, Title: "Set up SIEM alert for suspicious login patterns", AssigneeCSV: "Henry Walsh", Status: "Todo", Priority: "Medium"},
	})

	// ── Security Standups ─────────────────────────────────────────────────────
	db.Create(&[]models.StandupEntry{
		{ProjectID: security.ID, Member: "Henry Walsh", Role: "Security Eng", Yesterday: "Closed SQL injection finding with parameterised queries", Today: "Implement audit logging middleware for admin endpoints", Blockers: "None", Status: "On Track", Date: today},
		{ProjectID: security.ID, Member: "Alice Chen", Role: "Backend", Yesterday: "Analysed IDOR vulnerability scope", Today: "Add object-level authorisation checks to user API", Blockers: "None", Status: "On Track", Date: today},
		{ProjectID: security.ID, Member: "Dan Kim", Role: "DevOps", Yesterday: "Enforced MFA for all admin IAM accounts", Today: "Audit production secrets and schedule rotation", Blockers: "Vault upgrade needed before rotation — waiting on approval", Status: "At Risk", Date: today},
		{ProjectID: security.ID, Member: "Iris Novak", Role: "Compliance", Yesterday: "Mapped data flows for SOC 2 gap analysis", Today: "Write access control matrix documentation", Blockers: "None", Status: "On Track", Date: today},
		// yesterday
		{ProjectID: security.ID, Member: "Henry Walsh", Role: "Security Eng", Yesterday: "Reproduced SQL injection in local env", Today: "Fix with parameterised queries and add regression test", Blockers: "None", Status: "On Track", Date: yesterday},
		{ProjectID: security.ID, Member: "Iris Novak", Role: "Compliance", Yesterday: "Kickoff meeting with external SOC 2 auditor", Today: "Map data flows for gap analysis", Blockers: "None", Status: "On Track", Date: yesterday},
	})

	// ── Security Deadlines ────────────────────────────────────────────────────
	db.Create(&[]models.Deadline{
		{ProjectID: security.ID, Title: "Pen Test Remediation Complete", Project: "Security", DueDateRaw: "2026-04-25", Priority: "Critical"},
		{ProjectID: security.ID, Title: "SOC 2 Type II Readiness Review", Project: "Security", DueDateRaw: "2026-05-10", Priority: "High"},
		{ProjectID: security.ID, Title: "Annual Security Awareness Training", Project: "Security", DueDateRaw: "2026-05-30", Priority: "Medium"},
	})

	// ── Security Meetings ─────────────────────────────────────────────────────
	sm1 := &models.Meeting{ProjectID: security.ID, Title: "Pen Test Debrief", Date: "Apr 10, 2026", AttendeeCSV: "Henry Walsh,Alice Chen,Dan Kim,Iris Novak", Notes: "External pen test returned 3 critical findings: SQL injection in search API, IDOR in user endpoints, weak session fixation. Remediation plan assigned."}
	sm2 := &models.Meeting{ProjectID: security.ID, Title: "SOC 2 Auditor Kickoff", Date: "Apr 14, 2026", AttendeeCSV: "Iris Novak,Henry Walsh,Dan Kim", Notes: "External auditor (Schellman) scoped the Type II assessment. Evidence collection window: May–July. Key controls to document: access management, change management, availability."}
	for _, m := range []*models.Meeting{sm1, sm2} {
		db.Create(m)
	}
	db.Create(&models.ActionItem{MeetingID: sm1.ID, Task: "Fix SQL injection — parameterise search query", Owner: "Henry Walsh", DueDate: "Apr 14", Done: true})
	db.Create(&models.ActionItem{MeetingID: sm1.ID, Task: "Fix IDOR — add object-level auth to user API", Owner: "Alice Chen", DueDate: "Apr 18", Done: false})
	db.Create(&models.ActionItem{MeetingID: sm1.ID, Task: "Fix session fixation — regenerate session on login", Owner: "Alice Chen", DueDate: "Apr 16", Done: true})
	db.Create(&models.ActionItem{MeetingID: sm2.ID, Task: "Draft access control matrix document", Owner: "Iris Novak", DueDate: "Apr 21", Done: false})
	db.Create(&models.ActionItem{MeetingID: sm2.ID, Task: "Export existing audit log data for auditor", Owner: "Dan Kim", DueDate: "Apr 25", Done: false})

	// ── Security Dev Tasks ────────────────────────────────────────────────────
	db.Create(&[]models.DevTask{
		{ProjectID: security.ID, Title: "Evaluate HashiCorp Vault vs AWS Secrets Manager", Type: "Research", AssigneeCSV: "Dan Kim", Status: "In Progress", Priority: "High"},
		{ProjectID: security.ID, Title: "Add rate limiting to authentication endpoints", Type: "Improvement", AssigneeCSV: "Henry Walsh", Status: "Todo", Priority: "High"},
		{ProjectID: security.ID, Title: "Remove hardcoded credentials from legacy config files", Type: "Tech Debt", AssigneeCSV: "Dan Kim", Status: "Done", Priority: "Critical"},
		{ProjectID: security.ID, Title: "Implement RBAC for API gateway", Type: "Improvement", AssigneeCSV: "Alice Chen", Status: "In Progress", Priority: "High"},
	})

	// ── Security Releases ─────────────────────────────────────────────────────
	srel1 := &models.Release{ProjectID: security.ID, Name: "Sec Patch 2026-04 – Pen Test Fixes", Description: "Remediation for critical pen test findings", Status: "In Progress", TargetDate: "2026-04-25"}
	db.Create(srel1)
	ss1 := &models.ReleaseStage{ReleaseID: srel1.ID, Name: "Security Review", Status: "Active"}
	ss2 := &models.ReleaseStage{ReleaseID: srel1.ID, Name: "Staging Verification", Status: "Pending"}
	ss3 := &models.ReleaseStage{ReleaseID: srel1.ID, Name: "Production Deploy", Status: "Pending"}
	db.Create(ss1); db.Create(ss2); db.Create(ss3)
	db.Create(&models.ReleaseStory{StageID: ss1.ID, Title: "SQL injection fix — parameterised search query", Assignee: "Henry Walsh", Status: "Passed"})
	db.Create(&models.ReleaseStory{StageID: ss1.ID, Title: "Session fixation fix — regenerate session on login", Assignee: "Alice Chen", Status: "Passed"})
	db.Create(&models.ReleaseStory{StageID: ss1.ID, Title: "IDOR fix — object-level authorisation on user API", Assignee: "Alice Chen", Status: "In QA"})
	db.Create(&models.ReleaseSlackUpdate{StageID: ss1.ID, Channel: "#security-releases", Message: "SQL injection and session fixation fixes verified by Henry. IDOR fix in review.", Author: "Henry Walsh", PostedAt: "Apr 17, 3:00 PM"})

	// ── Global (cross-project) meeting ────────────────────────────────────────
	gm := &models.Meeting{ProjectID: 0, Title: "All-Hands Engineering Sync", Date: "Apr 15, 2026", AttendeeCSV: "Alice Chen,Bob Martinez,Carol Singh,Dan Kim,Eva Park,Frank Liu,Grace Obi,Henry Walsh,Iris Novak", Notes: "Quarterly engineering all-hands. Discussed Q2 roadmap, hiring plan (3 backend, 1 security), and office days policy update."}
	db.Create(gm)
	db.Create(&models.ActionItem{MeetingID: gm.ID, Task: "Share Q2 roadmap slides with team", Owner: "EM", DueDate: "Apr 16", Done: true})
	db.Create(&models.ActionItem{MeetingID: gm.ID, Task: "Update job descriptions and post openings", Owner: "EM", DueDate: "Apr 22", Done: false})

	return nil
}
