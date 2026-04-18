-- =============================================================================
-- Sprinto — Seed Data
-- Generated for: 2026-04-18
-- =============================================================================
-- Run with: psql -d sprinto -f seed.sql
-- Assumes tables already exist (run the Go app once to auto-migrate, or provide
-- your own CREATE TABLE statements). This script truncates all tables first.
-- =============================================================================

BEGIN;

-- ── Truncate in dependency order ──────────────────────────────────────────────
TRUNCATE TABLE
  release_slack_updates,
  release_stories,
  release_stages,
  releases,
  dev_tasks,
  meeting_action_items,
  meetings,
  deadlines,
  standups,
  sprint_tasks,
  sprints,
  project_members,
  team_members,
  projects
RESTART IDENTITY CASCADE;

-- =============================================================================
-- TEAM MEMBERS
-- =============================================================================

INSERT INTO team_members (id, name, role, email, created_at, updated_at) VALUES
  (1,  'Alice Chen',    'Backend Engineer',     'alice@example.com',  NOW(), NOW()),
  (2,  'Bob Martinez',  'Frontend Engineer',    'bob@example.com',    NOW(), NOW()),
  (3,  'Carol Singh',   'Full Stack Engineer',  'carol@example.com',  NOW(), NOW()),
  (4,  'Dan Kim',       'DevOps Engineer',      'dan@example.com',    NOW(), NOW()),
  (5,  'Eva Park',      'QA Engineer',          'eva@example.com',    NOW(), NOW()),
  (6,  'Frank Liu',     'iOS Engineer',         'frank@example.com',  NOW(), NOW()),
  (7,  'Grace Obi',     'Android Engineer',     'grace@example.com',  NOW(), NOW()),
  (8,  'Henry Walsh',   'Security Engineer',    'henry@example.com',  NOW(), NOW()),
  (9,  'Iris Novak',    'Compliance Analyst',   'iris@example.com',   NOW(), NOW());

-- =============================================================================
-- PROJECTS
-- =============================================================================

INSERT INTO projects (id, name, description, created_at, updated_at) VALUES
  (1, 'Platform', 'Core API and backend infrastructure services', NOW(), NOW()),
  (2, 'Mobile',   'iOS and Android applications',                  NOW(), NOW()),
  (3, 'Security', 'Auth, compliance, and security hardening',      NOW(), NOW());

-- ── Project ↔ Member join table ───────────────────────────────────────────────
-- Platform: Alice(1), Bob(2), Carol(3), Dan(4), Eva(5)
-- Mobile:   Bob(2), Carol(3), Frank(6), Grace(7), Eva(5)
-- Security: Alice(1), Dan(4), Henry(8), Iris(9)

INSERT INTO project_members (project_id, team_member_id) VALUES
  (1, 1), (1, 2), (1, 3), (1, 4), (1, 5),
  (2, 2), (2, 3), (2, 6), (2, 7), (2, 5),
  (3, 1), (3, 4), (3, 8), (3, 9);

-- =============================================================================
-- SPRINTS
-- =============================================================================

INSERT INTO sprints (id, project_id, name, goal, progress, start_date, end_date, active, created_at, updated_at) VALUES
  (1, 1, 'Sprint 12', 'Ship OAuth2 integration and complete API v2 endpoints',               62, 'Apr 14', 'Apr 28', TRUE,  NOW(), NOW()),
  (2, 2, 'Sprint 8',  'Launch push notifications and implement offline mode for core screens', 41, 'Apr 14', 'Apr 28', TRUE,  NOW(), NOW()),
  (3, 3, 'Sprint 5',  'Complete penetration test remediation and achieve SOC 2 Type II readiness', 28, 'Apr 7', 'Apr 25', TRUE, NOW(), NOW());

-- =============================================================================
-- SPRINT TASKS
-- =============================================================================

-- ── Platform — Sprint 12 ──────────────────────────────────────────────────────
INSERT INTO sprint_tasks (sprint_id, title, assignee, status, priority, created_at, updated_at) VALUES
  (1, 'Fix authentication token refresh',    'Alice Chen',   'In Progress', 'High',   NOW(), NOW()),
  (1, 'Implement OAuth2 provider',           'Bob Martinez', 'In Progress', 'High',   NOW(), NOW()),
  (1, 'Database migration script',           'Carol Singh',  'Done',        'Medium', NOW(), NOW()),
  (1, 'CI/CD pipeline optimisation',         'Dan Kim',      'Todo',        'Medium', NOW(), NOW()),
  (1, 'Write E2E test suite',                'Eva Park',     'Todo',        'High',   NOW(), NOW()),
  (1, 'API rate limiting',                   'Alice Chen',   'Done',        'Medium', NOW(), NOW()),
  (1, 'Mobile responsive fixes',             'Bob Martinez', 'Todo',        'Low',    NOW(), NOW()),
  (1, 'Cache layer implementation',          'Carol Singh',  'In Progress', 'High',   NOW(), NOW()),
  (1, 'Docker Compose setup',                'Dan Kim',      'Done',        'Low',    NOW(), NOW()),
  (1, 'Performance regression tests',        'Eva Park',     'In Progress', 'Medium', NOW(), NOW());

-- ── Mobile — Sprint 8 ─────────────────────────────────────────────────────────
INSERT INTO sprint_tasks (sprint_id, title, assignee, status, priority, created_at, updated_at) VALUES
  (2, 'Push notification service integration',  'Frank Liu',    'In Progress', 'High',   NOW(), NOW()),
  (2, 'Offline data sync for home screen',      'Grace Obi',    'In Progress', 'High',   NOW(), NOW()),
  (2, 'iOS deep link handling',                 'Frank Liu',    'Todo',        'Medium', NOW(), NOW()),
  (2, 'Android background fetch',               'Grace Obi',    'Todo',        'Medium', NOW(), NOW()),
  (2, 'Notification permission flow UI',        'Bob Martinez', 'Done',        'High',   NOW(), NOW()),
  (2, 'SQLite offline schema migration',        'Carol Singh',  'In Progress', 'High',   NOW(), NOW()),
  (2, 'E2E tests for notification scenarios',   'Eva Park',     'Todo',        'Medium', NOW(), NOW()),
  (2, 'App Store release notes draft',          'Bob Martinez', 'Todo',        'Low',    NOW(), NOW());

-- ── Security — Sprint 5 ───────────────────────────────────────────────────────
INSERT INTO sprint_tasks (sprint_id, title, assignee, status, priority, created_at, updated_at) VALUES
  (3, 'Remediate pen test finding: SQL injection in search',  'Henry Walsh', 'Done',        'Critical', NOW(), NOW()),
  (3, 'Remediate pen test finding: IDOR in user API',         'Alice Chen',  'In Progress', 'Critical', NOW(), NOW()),
  (3, 'Implement audit log for all admin actions',            'Henry Walsh', 'In Progress', 'High',     NOW(), NOW()),
  (3, 'Enable MFA enforcement for admin accounts',            'Dan Kim',     'Done',        'High',     NOW(), NOW()),
  (3, 'Review and rotate all production secrets',             'Dan Kim',     'Todo',        'High',     NOW(), NOW()),
  (3, 'Document access control matrix for SOC 2',            'Iris Novak',  'In Progress', 'Medium',   NOW(), NOW()),
  (3, 'Set up SIEM alert for suspicious login patterns',      'Henry Walsh', 'Todo',        'Medium',   NOW(), NOW());

-- =============================================================================
-- STANDUPS  (table: standups)
-- =============================================================================

-- ── Platform — today (2026-04-18) ────────────────────────────────────────────
INSERT INTO standups (project_id, member, role, yesterday, today, blockers, status, date, created_at, updated_at) VALUES
  (1, 'Alice Chen',   'Backend',   'Fixed JWT expiry bug in auth middleware',                'Continue OAuth2 integration with provider SDK',         'None',                                        'On Track', '2026-04-18', NOW(), NOW()),
  (1, 'Bob Martinez', 'Frontend',  'Updated UI component library to v3',                    'Implement dashboard charts and sprint progress view',    'Waiting on design approval for new layout',   'At Risk',  '2026-04-18', NOW(), NOW()),
  (1, 'Carol Singh',  'Full Stack','Completed DB migration script and tested in staging',   'Start cache layer implementation with Redis',            'None',                                        'On Track', '2026-04-18', NOW(), NOW()),
  (1, 'Dan Kim',      'DevOps',    'Reviewed Kubernetes config and updated resource limits', 'Optimise CI/CD pipeline run times',                     'Need access to production cluster',           'Blocked',  '2026-04-18', NOW(), NOW()),
  (1, 'Eva Park',     'QA',        'Created test plan for Sprint 12 features',               'Set up Playwright E2E test framework',                   'None',                                        'On Track', '2026-04-18', NOW(), NOW());

-- ── Platform — yesterday (2026-04-17) ────────────────────────────────────────
INSERT INTO standups (project_id, member, role, yesterday, today, blockers, status, date, created_at, updated_at) VALUES
  (1, 'Alice Chen',   'Backend',   'Reviewed OAuth2 spec and opened provider account',  'Start token refresh fix',              'None', 'On Track', '2026-04-17', NOW(), NOW()),
  (1, 'Bob Martinez', 'Frontend',  'Audited existing UI components',                    'Upgrade component library',            'None', 'On Track', '2026-04-17', NOW(), NOW()),
  (1, 'Carol Singh',  'Full Stack','Wrote DB migration plan',                           'Implement and test migration script',  'None', 'On Track', '2026-04-17', NOW(), NOW());

-- ── Platform — two days ago (2026-04-16) ─────────────────────────────────────
INSERT INTO standups (project_id, member, role, yesterday, today, blockers, status, date, created_at, updated_at) VALUES
  (1, 'Alice Chen', 'Backend', 'Sprint planning and task breakdown',   'Begin OAuth2 research',        'None', 'On Track', '2026-04-16', NOW(), NOW()),
  (1, 'Dan Kim',    'DevOps',  'Upgraded staging K8s cluster',         'Review prod cluster config',   'None', 'On Track', '2026-04-16', NOW(), NOW());

-- ── Mobile — today (2026-04-18) ──────────────────────────────────────────────
INSERT INTO standups (project_id, member, role, yesterday, today, blockers, status, date, created_at, updated_at) VALUES
  (2, 'Frank Liu',    'iOS',        'Integrated APNs token registration',                'Handle notification payload parsing and display',           'None',                                        'On Track', '2026-04-18', NOW(), NOW()),
  (2, 'Grace Obi',    'Android',    'Set up WorkManager for background sync',             'Implement conflict resolution for offline data sync',       'FCM quota limit hit in dev — waiting for increase', 'At Risk',  '2026-04-18', NOW(), NOW()),
  (2, 'Bob Martinez', 'Frontend',   'Shipped notification permission UI screens',         'Polish onboarding flow animations',                         'None',                                        'On Track', '2026-04-18', NOW(), NOW()),
  (2, 'Carol Singh',  'Full Stack', 'Wrote SQLite schema for offline tables',             'Run migration on device simulators and fix edge cases',     'None',                                        'On Track', '2026-04-18', NOW(), NOW()),
  (2, 'Eva Park',     'QA',         'Explored notification test tooling',                 'Write test cases for notification permission scenarios',     'None',                                        'On Track', '2026-04-18', NOW(), NOW());

-- ── Mobile — yesterday (2026-04-17) ──────────────────────────────────────────
INSERT INTO standups (project_id, member, role, yesterday, today, blockers, status, date, created_at, updated_at) VALUES
  (2, 'Frank Liu', 'iOS',     'Set up APNs certificates',         'Integrate APNs token registration',   'None', 'On Track', '2026-04-17', NOW(), NOW()),
  (2, 'Grace Obi', 'Android', 'Researched WorkManager API',       'Set up WorkManager for background sync', 'None', 'On Track', '2026-04-17', NOW(), NOW());

-- ── Security — today (2026-04-18) ────────────────────────────────────────────
INSERT INTO standups (project_id, member, role, yesterday, today, blockers, status, date, created_at, updated_at) VALUES
  (3, 'Henry Walsh', 'Security Eng', 'Closed SQL injection finding with parameterised queries',  'Implement audit logging middleware for admin endpoints',    'None',                                              'On Track', '2026-04-18', NOW(), NOW()),
  (3, 'Alice Chen',  'Backend',      'Analysed IDOR vulnerability scope',                        'Add object-level authorisation checks to user API',         'None',                                              'On Track', '2026-04-18', NOW(), NOW()),
  (3, 'Dan Kim',     'DevOps',       'Enforced MFA for all admin IAM accounts',                  'Audit production secrets and schedule rotation',            'Vault upgrade needed before rotation — waiting on approval', 'At Risk',  '2026-04-18', NOW(), NOW()),
  (3, 'Iris Novak',  'Compliance',   'Mapped data flows for SOC 2 gap analysis',                 'Write access control matrix documentation',                 'None',                                              'On Track', '2026-04-18', NOW(), NOW());

-- ── Security — yesterday (2026-04-17) ────────────────────────────────────────
INSERT INTO standups (project_id, member, role, yesterday, today, blockers, status, date, created_at, updated_at) VALUES
  (3, 'Henry Walsh', 'Security Eng', 'Reproduced SQL injection in local env',          'Fix with parameterised queries and add regression test',  'None', 'On Track', '2026-04-17', NOW(), NOW()),
  (3, 'Iris Novak',  'Compliance',   'Kickoff meeting with external SOC 2 auditor',    'Map data flows for gap analysis',                         'None', 'On Track', '2026-04-17', NOW(), NOW());

-- =============================================================================
-- DEADLINES
-- =============================================================================

-- ── Platform ──────────────────────────────────────────────────────────────────
INSERT INTO deadlines (project_id, title, project, due_date, priority, created_at, updated_at) VALUES
  (1, 'API v2 Public Launch',              'Platform', '2026-04-20', 'Critical', NOW(), NOW()),
  (1, 'OAuth2 Integration Complete',       'Platform', '2026-04-28', 'High',     NOW(), NOW()),
  (1, 'Q2 Engineering OKR Review',         'Platform', '2026-05-05', 'Medium',   NOW(), NOW()),
  (1, 'Infrastructure Upgrade — Phase 2',  'Platform', '2026-05-20', 'Low',      NOW(), NOW());

-- ── Mobile ────────────────────────────────────────────────────────────────────
INSERT INTO deadlines (project_id, title, project, due_date, priority, created_at, updated_at) VALUES
  (2, 'App Store Submission — v1.2',         'Mobile', '2026-04-25', 'Critical', NOW(), NOW()),
  (2, 'Beta TestFlight Release',             'Mobile', '2026-04-22', 'High',     NOW(), NOW()),
  (2, 'QA Sign-off for Push Notifications',  'Mobile', '2026-04-21', 'High',     NOW(), NOW()),
  (2, 'Google Play Internal Testing',        'Mobile', '2026-05-01', 'Medium',   NOW(), NOW());

-- ── Security ──────────────────────────────────────────────────────────────────
INSERT INTO deadlines (project_id, title, project, due_date, priority, created_at, updated_at) VALUES
  (3, 'Pen Test Remediation Complete',      'Security', '2026-04-25', 'Critical', NOW(), NOW()),
  (3, 'SOC 2 Type II Readiness Review',     'Security', '2026-05-10', 'High',     NOW(), NOW()),
  (3, 'Annual Security Awareness Training', 'Security', '2026-05-30', 'Medium',   NOW(), NOW());

-- =============================================================================
-- MEETINGS
-- =============================================================================

-- ── Platform meetings ─────────────────────────────────────────────────────────
INSERT INTO meetings (id, project_id, title, date, attendees, notes, created_at, updated_at) VALUES
  (1, 1, 'Sprint 12 Planning',
     'Apr 14, 2026',
     'Alice Chen,Bob Martinez,Carol Singh,Dan Kim,Eva Park',
     'Defined sprint goals focused on OAuth2 and API v2. Capacity at 85% due to Dan''s infrastructure work. Story points committed: 42.',
     NOW(), NOW()),
  (2, 1, '1:1 — Alice Chen',
     'Apr 15, 2026',
     'Alice Chen',
     'Discussed career growth path toward Staff Engineer. Alice interested in leading API v2 architecture initiative next quarter.',
     NOW(), NOW()),
  (3, 1, 'Architecture Review — API v2',
     'Apr 16, 2026',
     'Alice Chen,Bob Martinez,Carol Singh',
     'Decided to keep REST with improved versioning strategy and OpenAPI docs. GraphQL deferred to Q3.',
     NOW(), NOW()),
  (4, 1, 'Incident Post-mortem',
     'Apr 17, 2026',
     'Alice Chen,Bob Martinez,Carol Singh,Dan Kim,Eva Park',
     'Auth service outage Apr 16 (45 min). Root cause: TLS cert expired without alert. No data loss. SLA impacted.',
     NOW(), NOW());

-- ── Mobile meetings ───────────────────────────────────────────────────────────
INSERT INTO meetings (id, project_id, title, date, attendees, notes, created_at, updated_at) VALUES
  (5, 2, 'Sprint 8 Kickoff',
     'Apr 14, 2026',
     'Frank Liu,Grace Obi,Bob Martinez,Carol Singh,Eva Park',
     'Aligned on push notification architecture. Decided to use Firebase for Android and APNs for iOS with a shared backend abstraction layer.',
     NOW(), NOW()),
  (6, 2, 'App Store Review Prep',
     'Apr 16, 2026',
     'Frank Liu,Bob Martinez',
     'Reviewed App Store guidelines for notification features. Privacy manifest needs updating — required for submission.',
     NOW(), NOW());

-- ── Security meetings ─────────────────────────────────────────────────────────
INSERT INTO meetings (id, project_id, title, date, attendees, notes, created_at, updated_at) VALUES
  (7, 3, 'Pen Test Debrief',
     'Apr 10, 2026',
     'Henry Walsh,Alice Chen,Dan Kim,Iris Novak',
     'External pen test returned 3 critical findings: SQL injection in search API, IDOR in user endpoints, weak session fixation. Remediation plan assigned.',
     NOW(), NOW()),
  (8, 3, 'SOC 2 Auditor Kickoff',
     'Apr 14, 2026',
     'Iris Novak,Henry Walsh,Dan Kim',
     'External auditor (Schellman) scoped the Type II assessment. Evidence collection window: May–July. Key controls to document: access management, change management, availability.',
     NOW(), NOW());

-- ── Global meeting (visible to all projects — project_id = 0) ────────────────
INSERT INTO meetings (id, project_id, title, date, attendees, notes, created_at, updated_at) VALUES
  (9, 0, 'All-Hands Engineering Sync',
     'Apr 15, 2026',
     'Alice Chen,Bob Martinez,Carol Singh,Dan Kim,Eva Park,Frank Liu,Grace Obi,Henry Walsh,Iris Novak',
     'Quarterly engineering all-hands. Discussed Q2 roadmap, hiring plan (3 backend, 1 security), and office days policy update.',
     NOW(), NOW());

-- =============================================================================
-- MEETING ACTION ITEMS  (table: meeting_action_items)
-- =============================================================================

-- Platform — Sprint 12 Planning (meeting 1)
INSERT INTO meeting_action_items (meeting_id, task, owner, due_date, done, created_at, updated_at) VALUES
  (1, 'Share OAuth2 spec doc with team',            'Alice Chen',   'Apr 15', TRUE,  NOW(), NOW()),
  (1, 'Update design mockups for dashboard',        'Bob Martinez', 'Apr 16', FALSE, NOW(), NOW());

-- Platform — Architecture Review (meeting 3)
INSERT INTO meeting_action_items (meeting_id, task, owner, due_date, done, created_at, updated_at) VALUES
  (3, 'Write OpenAPI schema for v2 endpoints',      'Alice Chen',   'Apr 21', FALSE, NOW(), NOW());

-- Platform — Incident Post-mortem (meeting 4)
INSERT INTO meeting_action_items (meeting_id, task, owner, due_date, done, created_at, updated_at) VALUES
  (4, 'Set up cert expiry monitoring alerts',       'Dan Kim',      'Apr 20', FALSE, NOW(), NOW()),
  (4, 'Write incident summary for stakeholders',    'EM',           'Apr 18', FALSE, NOW(), NOW());

-- Mobile — Sprint 8 Kickoff (meeting 5)
INSERT INTO meeting_action_items (meeting_id, task, owner, due_date, done, created_at, updated_at) VALUES
  (5, 'Set up shared push notification abstraction in backend', 'Carol Singh', 'Apr 17', TRUE,  NOW(), NOW()),
  (5, 'Create Firebase project and share credentials',          'Grace Obi',   'Apr 15', TRUE,  NOW(), NOW());

-- Mobile — App Store Review Prep (meeting 6)
INSERT INTO meeting_action_items (meeting_id, task, owner, due_date, done, created_at, updated_at) VALUES
  (6, 'Update PrivacyInfo.xcprivacy manifest',      'Frank Liu',    'Apr 19', FALSE, NOW(), NOW()),
  (6, 'Screenshot all new notification UI screens', 'Bob Martinez', 'Apr 20', FALSE, NOW(), NOW());

-- Security — Pen Test Debrief (meeting 7)
INSERT INTO meeting_action_items (meeting_id, task, owner, due_date, done, created_at, updated_at) VALUES
  (7, 'Fix SQL injection — parameterise search query',            'Henry Walsh', 'Apr 14', TRUE,  NOW(), NOW()),
  (7, 'Fix IDOR — add object-level auth to user API',             'Alice Chen',  'Apr 18', FALSE, NOW(), NOW()),
  (7, 'Fix session fixation — regenerate session on login',       'Alice Chen',  'Apr 16', TRUE,  NOW(), NOW());

-- Security — SOC 2 Auditor Kickoff (meeting 8)
INSERT INTO meeting_action_items (meeting_id, task, owner, due_date, done, created_at, updated_at) VALUES
  (8, 'Draft access control matrix document',       'Iris Novak',   'Apr 21', FALSE, NOW(), NOW()),
  (8, 'Export existing audit log data for auditor', 'Dan Kim',      'Apr 25', FALSE, NOW(), NOW());

-- Global — All-Hands (meeting 9)
INSERT INTO meeting_action_items (meeting_id, task, owner, due_date, done, created_at, updated_at) VALUES
  (9, 'Share Q2 roadmap slides with team',          'EM',           'Apr 16', TRUE,  NOW(), NOW()),
  (9, 'Update job descriptions and post openings',  'EM',           'Apr 22', FALSE, NOW(), NOW());

-- =============================================================================
-- DEV TASKS  (table: dev_tasks)
-- =============================================================================

-- ── Platform ──────────────────────────────────────────────────────────────────
INSERT INTO dev_tasks (project_id, title, type, assignee, status, priority, created_at, updated_at) VALUES
  (1, 'Migrate REST endpoints to OpenAPI spec',      'Improvement', 'Alice Chen',   'In Progress', 'High',   NOW(), NOW()),
  (1, 'Upgrade Go version to 1.22',                  'Tech Debt',   'Dan Kim',      'In Progress', 'Medium', NOW(), NOW()),
  (1, 'Add Swagger UI for API documentation',        'Improvement', 'Bob Martinez', 'Todo',        'Medium', NOW(), NOW()),
  (1, 'Remove deprecated auth middleware',           'Tech Debt',   'Carol Singh',  'Done',        'High',   NOW(), NOW()),
  (1, 'Benchmark critical database queries',         'Research',    'Alice Chen',   'Todo',        'Low',    NOW(), NOW()),
  (1, 'Setup Sentry error monitoring',               'Improvement', 'Dan Kim',      'Todo',        'High',   NOW(), NOW()),
  (1, 'Consolidate app config management',           'Tech Debt',   'Carol Singh',  'Todo',        'Medium', NOW(), NOW()),
  (1, 'Load testing with k6',                        'Research',    'Eva Park',     'Todo',        'Medium', NOW(), NOW());

-- ── Mobile ────────────────────────────────────────────────────────────────────
INSERT INTO dev_tasks (project_id, title, type, assignee, status, priority, created_at, updated_at) VALUES
  (2, 'Upgrade React Native to 0.74',                        'Tech Debt',   'Bob Martinez', 'Todo',        'High',   NOW(), NOW()),
  (2, 'Implement biometric auth (Face ID / fingerprint)',    'Improvement', 'Frank Liu',    'Todo',        'Medium', NOW(), NOW()),
  (2, 'Benchmark SQLite vs Realm for offline storage',       'Research',    'Carol Singh',  'In Progress', 'Medium', NOW(), NOW()),
  (2, 'Remove legacy Bluetooth module',                      'Tech Debt',   'Grace Obi',    'Done',        'Low',    NOW(), NOW()),
  (2, 'Add crash analytics with Crashlytics',                'Improvement', 'Frank Liu',    'Todo',        'High',   NOW(), NOW());

-- ── Security ──────────────────────────────────────────────────────────────────
INSERT INTO dev_tasks (project_id, title, type, assignee, status, priority, created_at, updated_at) VALUES
  (3, 'Evaluate HashiCorp Vault vs AWS Secrets Manager', 'Research',    'Dan Kim',     'In Progress', 'High',     NOW(), NOW()),
  (3, 'Add rate limiting to authentication endpoints',   'Improvement', 'Henry Walsh', 'Todo',        'High',     NOW(), NOW()),
  (3, 'Remove hardcoded credentials from legacy config', 'Tech Debt',   'Dan Kim',     'Done',        'Critical', NOW(), NOW()),
  (3, 'Implement RBAC for API gateway',                  'Improvement', 'Alice Chen',  'In Progress', 'High',     NOW(), NOW());

-- =============================================================================
-- RELEASES
-- =============================================================================

INSERT INTO releases (id, project_id, name, description, status, target_date, created_at, updated_at) VALUES
  (1, 1, 'v2.3.0 – Auth Overhaul',         'OAuth2 integration and JWT improvements',                    'In Progress', '2026-04-28', NOW(), NOW()),
  (2, 1, 'v2.2.1 – Security Hotfix',       'TLS certificate rotation and auth session hardening',        'Released',    '2026-04-18', NOW(), NOW()),
  (3, 2, 'v1.2.0 – Push & Offline',        'Push notifications and offline mode for core screens',       'In Progress', '2026-04-25', NOW(), NOW()),
  (4, 3, 'Sec Patch 2026-04 – Pen Test Fixes', 'Remediation for critical pen test findings',             'In Progress', '2026-04-25', NOW(), NOW());

-- =============================================================================
-- RELEASE STAGES  (table: release_stages)
-- =============================================================================

INSERT INTO release_stages (id, release_id, name, status, created_at, updated_at) VALUES
  -- v2.3.0 Auth Overhaul (release 1)
  (1,  1, 'QA Round 1',           'Done',    NOW(), NOW()),
  (2,  1, 'QA Round 2',           'Active',  NOW(), NOW()),
  (3,  1, 'Staging Deploy',       'Pending', NOW(), NOW()),
  -- v2.2.1 Security Hotfix (release 2)
  (4,  2, 'QA Verification',      'Done',    NOW(), NOW()),
  (5,  2, 'Production Deploy',    'Done',    NOW(), NOW()),
  -- v1.2.0 Push & Offline (release 3)
  (6,  3, 'Internal QA',                  'Active',  NOW(), NOW()),
  (7,  3, 'Beta — TestFlight / Play',     'Pending', NOW(), NOW()),
  (8,  3, 'App Store / Play Store',       'Pending', NOW(), NOW()),
  -- Sec Patch (release 4)
  (9,  4, 'Security Review',      'Active',  NOW(), NOW()),
  (10, 4, 'Staging Verification', 'Pending', NOW(), NOW()),
  (11, 4, 'Production Deploy',    'Pending', NOW(), NOW());

-- =============================================================================
-- RELEASE STORIES  (table: release_stories)
-- =============================================================================

-- v2.3.0 — QA Round 1 (stage 1)
INSERT INTO release_stories (stage_id, title, assignee, status, created_at, updated_at) VALUES
  (1, 'Fix authentication token refresh', 'Eva Park', 'Passed', NOW(), NOW()),
  (1, 'API rate limiting',                'Eva Park', 'Passed', NOW(), NOW());

-- v2.3.0 — QA Round 2 (stage 2)
INSERT INTO release_stories (stage_id, title, assignee, status, created_at, updated_at) VALUES
  (2, 'Implement OAuth2 provider',    'Eva Park', 'In QA',   NOW(), NOW()),
  (2, 'Cache layer implementation',   'Eva Park', 'Pending', NOW(), NOW());

-- v2.2.1 — QA Verification (stage 4)
INSERT INTO release_stories (stage_id, title, assignee, status, created_at, updated_at) VALUES
  (4, 'TLS cert auto-rotation script',    'Dan Kim',     'Passed', NOW(), NOW()),
  (4, 'Session token TTL enforcement',    'Alice Chen',  'Passed', NOW(), NOW());

-- v1.2.0 — Internal QA (stage 6)
INSERT INTO release_stories (stage_id, title, assignee, status, created_at, updated_at) VALUES
  (6, 'Push notification permission flow',        'Eva Park', 'Passed', NOW(), NOW()),
  (6, 'Push notification delivery — iOS',         'Eva Park', 'In QA',  NOW(), NOW()),
  (6, 'Push notification delivery — Android',     'Eva Park', 'Pending',NOW(), NOW());

-- Sec Patch — Security Review (stage 9)
INSERT INTO release_stories (stage_id, title, assignee, status, created_at, updated_at) VALUES
  (9, 'SQL injection fix — parameterised search query',       'Henry Walsh', 'Passed', NOW(), NOW()),
  (9, 'Session fixation fix — regenerate session on login',   'Alice Chen',  'Passed', NOW(), NOW()),
  (9, 'IDOR fix — object-level authorisation on user API',    'Alice Chen',  'In QA',  NOW(), NOW());

-- =============================================================================
-- RELEASE SLACK UPDATES  (table: release_slack_updates)
-- =============================================================================

-- v2.3.0 — QA Round 1 (stage 1)
INSERT INTO release_slack_updates (stage_id, channel, message, author, posted_at, created_at, updated_at) VALUES
  (1, '#releases', 'QA Round 1 complete — all stories passed. Promoting to Round 2.', 'Eva Park', 'Apr 16, 2:30 PM', NOW(), NOW());

-- v2.3.0 — QA Round 2 (stage 2)
INSERT INTO release_slack_updates (stage_id, channel, message, author, posted_at, created_at, updated_at) VALUES
  (2, '#releases', 'Starting QA Round 2. OAuth2 story handed to Eva. Cache layer follows once OAuth2 passes.', 'Alice Chen', 'Apr 17, 10:15 AM', NOW(), NOW());

-- v2.2.1 — QA Verification (stage 4)
INSERT INTO release_slack_updates (stage_id, channel, message, author, posted_at, created_at, updated_at) VALUES
  (4, '#releases',  'Hotfix verified in staging. Both stories passed QA. Requesting prod deploy window.', 'Eva Park', 'Apr 17, 4:00 PM', NOW(), NOW());

-- v2.2.1 — Production Deploy (stage 5)
INSERT INTO release_slack_updates (stage_id, channel, message, author, posted_at, created_at, updated_at) VALUES
  (5, '#incidents', 'v2.2.1 deployed to prod. TLS rotation confirmed. Monitoring for 30 min.',  'Dan Kim', 'Apr 17, 6:45 PM', NOW(), NOW()),
  (5, '#releases',  'v2.2.1 stable — no issues. Release complete.',                              'EM',      'Apr 17, 7:20 PM', NOW(), NOW());

-- v1.2.0 — Internal QA (stage 6)
INSERT INTO release_slack_updates (stage_id, channel, message, author, posted_at, created_at, updated_at) VALUES
  (6, '#mobile-releases', 'Internal QA started. Permission flow passed. iOS push delivery testing in progress.', 'Eva Park', 'Apr 17, 9:00 AM', NOW(), NOW());

-- Sec Patch — Security Review (stage 9)
INSERT INTO release_slack_updates (stage_id, channel, message, author, posted_at, created_at, updated_at) VALUES
  (9, '#security-releases', 'SQL injection and session fixation fixes verified by Henry. IDOR fix in review.', 'Henry Walsh', 'Apr 17, 3:00 PM', NOW(), NOW());

-- =============================================================================
-- Reset sequences so future INSERTs (without explicit IDs) work correctly
-- =============================================================================

SELECT setval('team_members_id_seq',           (SELECT MAX(id) FROM team_members));
SELECT setval('projects_id_seq',               (SELECT MAX(id) FROM projects));
SELECT setval('sprints_id_seq',                (SELECT MAX(id) FROM sprints));
SELECT setval('sprint_tasks_id_seq',           (SELECT MAX(id) FROM sprint_tasks));
SELECT setval('standups_id_seq',               (SELECT MAX(id) FROM standups));
SELECT setval('deadlines_id_seq',              (SELECT MAX(id) FROM deadlines));
SELECT setval('meetings_id_seq',               (SELECT MAX(id) FROM meetings));
SELECT setval('meeting_action_items_id_seq',   (SELECT MAX(id) FROM meeting_action_items));
SELECT setval('dev_tasks_id_seq',              (SELECT MAX(id) FROM dev_tasks));
SELECT setval('releases_id_seq',               (SELECT MAX(id) FROM releases));
SELECT setval('release_stages_id_seq',         (SELECT MAX(id) FROM release_stages));
SELECT setval('release_stories_id_seq',        (SELECT MAX(id) FROM release_stories));
SELECT setval('release_slack_updates_id_seq',  (SELECT MAX(id) FROM release_slack_updates));

COMMIT;
