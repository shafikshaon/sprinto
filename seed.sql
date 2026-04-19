-- ─────────────────────────────────────────────────────────────────────────────
-- Seed: Daily Standups  (19 Apr 2026 → 30 May 2026)
-- Run:  psql $DATABASE_URL -f seed.sql
-- Requires at least one row in users and projects tables.
-- Each standup is randomly assigned to one of the existing projects.
-- ─────────────────────────────────────────────────────────────────────────────

DO $$
DECLARE
  v_user_name   text;
  v_project_ids integer[];
  v_pid         integer;
  v_nproj       integer;

BEGIN
  -- ── 1. First user ───────────────────────────────────────────────────────────
  SELECT full_name
    INTO v_user_name
    FROM users
   WHERE deleted_at IS NULL
   ORDER BY id
   LIMIT 1;

  IF v_user_name IS NULL THEN
    RAISE EXCEPTION 'No users found — create at least one user before seeding.';
  END IF;

  -- ── 2. All project IDs ──────────────────────────────────────────────────────
  SELECT ARRAY(
    SELECT id FROM projects WHERE deleted_at IS NULL ORDER BY id
  ) INTO v_project_ids;

  v_nproj := array_length(v_project_ids, 1);
  IF v_nproj IS NULL THEN
    RAISE EXCEPTION 'No projects found — create at least one project before seeding.';
  END IF;

  RAISE NOTICE 'Seeding standups for user "%" across % project(s)…', v_user_name, v_nproj;

  -- ── 3. Standup entries ──────────────────────────────────────────────────────
  -- Weekdays 20 Apr → 29 May 2026 (skip weekends)

  -- Week 1 (20–24 Apr) — Sprint kick-off, auth module
  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-04-20',
   'Sprint 12 kick-off. Team aligned on goals: complete auth module, wire up new dashboard layout, and close 4 carry-over tickets from Sprint 11. Story-point split reviewed and accepted.',
   'Design team to deliver final Figma specs for dashboard by Wednesday',
   '',
   v_user_name || ' to create Sprint 12 board and assign tickets. Nadia to start JWT refresh-token flow.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-04-21',
   'JWT refresh-token endpoint completed. Arif started role-based middleware. Decided to use httpOnly cookie over header for token storage after brief discussion.',
   '',
   'Redis session store unavailable on staging — blocking session-invalidation tests.',
   'DevOps (Karim) to provision Redis on staging by EOD. Arif to stub session invalidation in the meantime.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-04-22',
   'Redis live on staging. Session invalidation tested and working. Role middleware merged. ' || v_user_name || ' reviewed Figma specs — minor layout questions raised with design.',
   'Awaiting design clarification on mobile breakpoint for sidebar',
   '',
   'Nadia to write integration tests for auth flow. ' || v_user_name || ' to reply to design by noon.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-04-23',
   'Auth integration tests all green. Password-reset email flow merged to main. Dashboard scaffold committed — top nav and sidebar responsive. PR review queue has 3 open items.',
   '',
   'Email delivery failing on staging (wrong SMTP credentials in .env.staging).',
   'Karim to update staging SMTP credentials. Team to clear PR review queue before EOD.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-04-24',
   'SMTP fixed — password reset emails working end-to-end. All 3 PRs reviewed and merged. Sprint 12 board is 28% complete. Auth module retrospective: went smoother than expected.',
   '',
   '',
   'Nadia starts dashboard data-fetching layer on Monday. Arif picks up notification service ticket.');

  -- Week 2 (27 Apr – 1 May) — Dashboard data layer
  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-04-27',
   'Nadia started dashboard API endpoints. Arif flagged notification service ticket is larger than estimated — needs re-pointing. ' || v_user_name || ' conducted 1-on-1s; no blockers raised.',
   'Backend API contract for dashboard widgets needs product sign-off',
   '',
   v_user_name || ' to get product sign-off on API contract today. Arif to split notification ticket into two.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-04-28',
   'API contract approved. 4 of 7 dashboard widget endpoints done. Arif split notification ticket — in-app notifications in scope, email digest moved to backlog. Sprint now 41% done.',
   '',
   'N+1 query issue on sprint-summary endpoint — tasks being fetched one by one.',
   'Nadia to add Preload("Tasks") to sprint-summary query. Arif to start in-app notification model.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-04-29',
   'N+1 fixed — sprint-summary endpoint now 12× faster. All 7 dashboard endpoints complete. Front-end (Rina) starting widget components. Notification model schema reviewed and merged.',
   'Front-end needs live staging data to verify widget rendering',
   '',
   'Karim to deploy latest backend to staging. Rina to wire widgets to staging API tomorrow.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-04-30',
   'Staging deployment done. Rina connected 3 of 5 dashboard widgets to live API — charts rendering correctly. Arif finished in-app notification dispatch logic. Sprint 58% complete.',
   '',
   'Chart.js date-axis timezone bug causing off-by-one on the daily standup chart.',
   'Rina to investigate Chart.js timezone config. Workaround: force UTC in dataset labels for now.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-01',
   'Chart.js timezone workaround shipped. All 5 dashboard widgets live on staging. Notification bell UI added by Rina. Sprint at 67% going into the long weekend.',
   '',
   '',
   'Team to do exploratory testing over the weekend if available. Sprint review call booked for Monday 10 AM.');

  -- Week 3 (4–8 May) — Testing & polish
  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-04',
   'Sprint review done — product happy with dashboard and notifications. 6 minor UI bugs filed from exploratory testing. 4 prioritised for this sprint, 2 deferred. Velocity: 62 points.',
   '',
   '',
   'Nadia takes bugs #47 and #48. Rina takes #49 and #51. Arif continues on notification read-receipts.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-05',
   'Bugs #47 (wrong task count) and #49 (sidebar flicker) fixed and in review. Arif added read-receipt tracking to notification model. Rina polishing empty-state illustrations.',
   'Waiting on copywriter for empty-state microcopy (3 screens)',
   '',
   v_user_name || ' to chase copywriter. Nadia to start bug #48 (date filter not persisting).');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-06',
   'Copywriter delivered microcopy. All 4 bugs resolved and merged. Notification read-receipts working end-to-end. Sprint at 89%. QA pass scheduled for tomorrow.',
   '',
   'Notification badge count not decrementing in real-time — needs polling or WebSocket.',
   'Decided on 30-second polling for now; WebSocket scoped for Sprint 13. Arif to implement polling today.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-07',
   'QA pass done — 2 low-severity findings, both fixed same day. Polling for notification badge merged. Sprint 97% done. Retrospective held: team praised async PR reviews, wants shorter standups.',
   '',
   '',
   v_user_name || ' to reduce standup timebox to 10 min from Sprint 13. Prep Sprint 13 planning for Friday.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-08',
   'Sprint 12 closed at 100%. Sprint 13 planning complete: release pipeline, CI/CD improvements, and user settings page. 54 points committed. Team energised after a successful sprint.',
   '',
   '',
   'Arif to set up GitHub Actions workflow for staging auto-deploy. Nadia to start user settings schema.');

  -- Week 4 (11–15 May) — CI/CD + user settings
  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-11',
   'GitHub Actions pipeline running — builds and tests on every PR. Nadia drafted user settings schema: profile, password change, notification prefs, timezone. Sprint 13 at 18%.',
   'Need DevOps approval for new IAM role for the Actions runner',
   '',
   'Karim to approve IAM role. Nadia to create migration file for settings schema.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-12',
   'IAM role approved. CI pipeline now deploys to staging on merge to main. Settings schema migration merged. Rina started settings page UI.',
   '',
   'Staging deploy takes 4 min — Docker image rebuild on every deploy is slow.',
   'Arif to add layer caching to Dockerfile. Expected to cut deploy time to ~90 seconds.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-13',
   'Docker layer caching done — staging deploy now 1m 40s. Profile and password-change endpoints complete. Rina finished settings page layout. Sprint 13 at 52%.',
   '',
   '',
   'Nadia to implement notification preferences endpoint. Rina to wire settings form to API.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-14',
   'Notification preferences endpoint done. Settings form working — profile save and password change functional on staging. Timezone preference saves but display not yet updated globally.',
   'Timezone display requires refactoring date helpers across 11 templates',
   '',
   v_user_name || ' to scope timezone refactor — may push to Sprint 14. Smoke test settings page tomorrow.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-15',
   'Smoke test passed with one issue: avatar upload silently failing on files > 2 MB. Fixed with proper validation and error message. Timezone refactor moved to Sprint 14. Sprint 13 at 81%.',
   '',
   '',
   'Remaining: CI prod pipeline, accessibility audit on settings page. Target done by Wednesday next week.');

  -- Week 5 (18–22 May) — Prod pipeline + v1.4 release
  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-18',
   'Arif building prod deployment pipeline — staging works, prod needs separate secrets and approval gate. Accessibility audit: 3 findings (missing labels, broken focus traps). Sprint 87%.',
   'Product sign-off needed before enabling prod auto-deploy',
   '',
   'Rina to fix accessibility findings. Arif to add manual approval gate to prod pipeline.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-19',
   'All 3 accessibility issues resolved. Prod pipeline approval gate tested with dummy release. Product sign-off received. Sprint 13 effectively done. Planning early release of v1.4.',
   '',
   '',
   v_user_name || ' to draft v1.4 release notes. Karim to do final infra check on prod. Release target: 22 May.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-20',
   'Release notes drafted and reviewed. Karim increased prod DB connection pool from 10 to 25. Full regression run — all 312 tests passing. RC tagged as v1.4.0-rc1.',
   '',
   'One flaky test in notification suite — intermittent timing issue, not blocking release.',
   'Nadia to add retry logic to flaky test after release. Go/no-go call tomorrow 9 AM.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-21',
   'Go/no-go: green across the board. Release scheduled tomorrow 10 AM. Rollback plan documented. On-call rotation set. Monitoring dashboards reviewed by the full team.',
   '',
   '',
   'All hands tomorrow 10 AM for release. Karim to watch infra metrics. ' || v_user_name || ' to send release comms.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-22',
   'v1.4.0 released to production at 10:17 AM. Zero-downtime deploy. All health checks green. Initial user feedback positive — settings page and dashboard well received.',
   '',
   'One user reported notification badge not updating after first login — root cause identified.',
   'Arif to patch notification badge init bug in hotfix. Target hotfix deploy today 4 PM.');

  -- Week 6 (25–29 May) — Post-release stabilisation
  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-25',
   'v1.4.1 hotfix deployed Friday 4 PM — notification badge resolved. Error rate back to baseline (<0.1%). Sprint 14 planning: timezone refactor, audit log, performance profiling. 48 points committed.',
   '',
   '',
   'Nadia leads timezone refactor. Arif starts audit log schema. Rina picks up performance profiling setup.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-26',
   'Timezone refactor: 6 of 11 templates updated. Arif completed audit log model and migration — tracks user, action, entity, and timestamp. Rina added pprof endpoint behind admin flag.',
   'UX decision needed on timezone display format (relative vs absolute)',
   '',
   v_user_name || ' to schedule UX call today. Nadia to continue templates once format is confirmed.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-27',
   'UX decision: show absolute local time with UTC offset in tooltips. All 11 templates updated — timezone now respects user preference globally. Rina ran first pprof profile; found slow query in audit log list.',
   '',
   'Audit log list query full-scanning table — missing index on (user_id, created_at).',
   'Arif to add composite index in new migration. Expected to drop query time from 420ms to ~15ms.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-28',
   'Composite index added — audit log query now 18ms. Timezone refactor fully merged and tested. Rina profiled 3 more endpoints: all within acceptable thresholds. Sprint 14 at 61%.',
   '',
   '',
   'Nadia to start audit log UI (admin panel). Arif to write service-layer tests for audit log.');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO standups (created_at, updated_at, project_id, date, summary, dependencies, issues, action_items) VALUES
  (NOW(), NOW(), v_pid, '2026-05-29',
   'Audit log UI draft in review — filterable by user and action type. All audit log service tests green. Performance profiling report shared with ' || v_user_name || '. Sprint 14 at 74%, on track.',
   '',
   '',
   'Continue audit log UI review. Mid-sprint check-in scheduled for Tuesday next week.');

  RAISE NOTICE 'Done — 30 standup entries inserted.';
END $$;

-- ─────────────────────────────────────────────────────────────────────────────
-- Seed: Dev Tasks
-- Randomly assigned to existing projects; assignees pulled from team_members.
-- ─────────────────────────────────────────────────────────────────────────────

DO $$
DECLARE
  v_project_ids  integer[];
  v_member_names text[];
  v_nproj        integer;
  v_nmem         integer;
  v_pid          integer;
  v_m1           text;
  v_m2           text;

BEGIN
  -- ── Projects ──────────────────────────────────────────────────────────────
  SELECT ARRAY(SELECT id FROM projects WHERE deleted_at IS NULL ORDER BY id)
    INTO v_project_ids;
  v_nproj := array_length(v_project_ids, 1);

  IF v_nproj IS NULL THEN
    RAISE EXCEPTION 'No projects found.';
  END IF;

  -- ── Team members (names only) ─────────────────────────────────────────────
  SELECT ARRAY(SELECT name FROM team_members WHERE deleted_at IS NULL ORDER BY id)
    INTO v_member_names;
  v_nmem := array_length(v_member_names, 1);

  IF v_nmem IS NULL THEN
    v_member_names := ARRAY['Unassigned'];
    v_nmem := 1;
  END IF;

  RAISE NOTICE 'Seeding dev tasks across % project(s) with % member(s)…', v_nproj, v_nmem;

  -- ── Improvement tasks ─────────────────────────────────────────────────────
  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_m2 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'38 days', NOW()-interval'10 days', v_pid, 'Migrate REST endpoints to GraphQL', 'Improvement', v_m1||','||v_m2, 'In Progress', 'High');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_m2 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'35 days', NOW()-interval'35 days', v_pid, 'Upgrade dashboard charts to Chart.js v4', 'Improvement', v_m1, 'Todo', 'Medium');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_m2 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'33 days', NOW()-interval'2 days', v_pid, 'Add dark mode support across all pages', 'Improvement', v_m1||','||v_m2, 'In Progress', 'Medium');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'30 days', NOW()-interval'30 days', v_pid, 'Implement keyboard shortcut navigation', 'Improvement', v_m1, 'Todo', 'Low');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_m2 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'28 days', NOW()-interval'5 days', v_pid, 'Optimise sprint board rendering performance', 'Improvement', v_m1||','||v_m2, 'Done', 'High');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'25 days', NOW()-interval'25 days', v_pid, 'Add CSV export for deadline reports', 'Improvement', v_m1, 'Todo', 'Medium');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_m2 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'22 days', NOW()-interval'8 days', v_pid, 'Redesign onboarding flow for new users', 'Improvement', v_m1||','||v_m2, 'In Progress', 'High');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'20 days', NOW()-interval'20 days', v_pid, 'Add real-time notification badge via WebSocket', 'Improvement', v_m1, 'Todo', 'High');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_m2 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'18 days', NOW()-interval'3 days', v_pid, 'Improve pagination UX with infinite scroll option', 'Improvement', v_m1||','||v_m2, 'Done', 'Medium');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'15 days', NOW()-interval'15 days', v_pid, 'Add multi-language (i18n) support skeleton', 'Improvement', v_m1, 'Todo', 'Low');

  -- ── Tech Debt tasks ───────────────────────────────────────────────────────
  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'40 days', NOW()-interval'40 days', v_pid, 'Remove deprecated jQuery dependency from legacy pages', 'Tech Debt', v_m1, 'Todo', 'High');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_m2 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'37 days', NOW()-interval'12 days', v_pid, 'Consolidate duplicate date-formatting helpers across templates', 'Tech Debt', v_m1||','||v_m2, 'In Progress', 'Medium');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'34 days', NOW()-interval'4 days', v_pid, 'Replace raw SQL queries in reporting module with GORM', 'Tech Debt', v_m1, 'Done', 'Critical');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'31 days', NOW()-interval'31 days', v_pid, 'Fix N+1 query on team members listing page', 'Tech Debt', v_m1, 'Todo', 'High');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_m2 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'27 days', NOW()-interval'9 days', v_pid, 'Add missing database indexes on foreign key columns', 'Tech Debt', v_m1||','||v_m2, 'Done', 'Critical');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'24 days', NOW()-interval'24 days', v_pid, 'Migrate hardcoded config values to environment variables', 'Tech Debt', v_m1, 'Todo', 'Medium');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_m2 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'21 days', NOW()-interval'6 days', v_pid, 'Standardise HTTP error response format across all handlers', 'Tech Debt', v_m1||','||v_m2, 'In Progress', 'Medium');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'17 days', NOW()-interval'17 days', v_pid, 'Remove unused feature-flag code from Sprint 9', 'Tech Debt', v_m1, 'Todo', 'Low');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_m2 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'14 days', NOW()-interval'14 days', v_pid, 'Increase unit test coverage for auth service to 80%', 'Tech Debt', v_m1||','||v_m2, 'Todo', 'High');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'10 days', NOW()-interval'2 days', v_pid, 'Upgrade Go module dependencies to latest stable versions', 'Tech Debt', v_m1, 'In Progress', 'Medium');

  -- ── Research tasks ────────────────────────────────────────────────────────
  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'36 days', NOW()-interval'36 days', v_pid, 'Evaluate OpenTelemetry for distributed tracing', 'Research', v_m1, 'Todo', 'Medium');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_m2 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'32 days', NOW()-interval'7 days', v_pid, 'Proof of concept: edge caching with Cloudflare Workers', 'Research', v_m1||','||v_m2, 'In Progress', 'High');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'29 days', NOW()-interval'29 days', v_pid, 'Assess feasibility of WASM for client-side report rendering', 'Research', v_m1, 'Todo', 'Low');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_m2 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'26 days', NOW()-interval'11 days', v_pid, 'Benchmark PostgreSQL JSONB vs relational for audit log', 'Research', v_m1||','||v_m2, 'Done', 'Medium');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'23 days', NOW()-interval'23 days', v_pid, 'Investigate AI-assisted sprint velocity prediction', 'Research', v_m1, 'Todo', 'Low');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_m2 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'19 days', NOW()-interval'5 days', v_pid, 'Research headless CMS options for documentation portal', 'Research', v_m1||','||v_m2, 'In Progress', 'Medium');

  v_m1 := v_member_names[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW()-interval'16 days', NOW()-interval'16 days', v_pid, 'Explore event-driven architecture with NATS for notifications', 'Research', v_m1, 'Todo', 'High');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW() - interval '12 days', NOW() - interval '12 days', v_pid,
   'Compare ClickHouse vs TimescaleDB for analytics pipeline',
   'Research', v_m1||','||v_m2, 'Todo', 'Medium');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW() - interval '8 days', NOW() - interval '1 day', v_pid,
   'Evaluate Playwright vs Cypress for end-to-end test suite',
   'Research', v_m1, 'In Progress', 'Medium');

  v_pid := v_project_ids[1 + floor(random() * v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, assignees, status, priority) VALUES
  (NOW() - interval '5 days', NOW() - interval '5 days', v_pid,
   'Spike: feature flags service (LaunchDarkly vs homegrown)',
   'Research', v_m1||','||v_m2, 'Todo', 'Low');

  RAISE NOTICE 'Done — 30 dev task entries inserted.';
END $$;

-- ─────────────────────────────────────────────────────────────────────────────
-- Seed: Dev Task Comments (standalone — safe to run on an existing DB)
-- ─────────────────────────────────────────────────────────────────────────────

DO $$
DECLARE
  v_task_ids     integer[];
  v_member_names text[];
  v_nmem         integer;
  v_ntasks       integer;
  v_m1           text;
  v_m2           text;
  v_m3           text;
  v_tid          integer;

  v_comments_a   text[] := ARRAY[
    'Started looking into this. Will share findings by end of week.',
    'Picked this up. Dependencies look manageable — nothing blocking yet.',
    'Scoped this out. Estimates updated in the sprint board.',
    'Quick spike done. Approach is feasible, moving to implementation.',
    'Reviewed the existing code. A few things need cleanup before we can start properly.',
    'Drafted an initial plan. Would appreciate a review from anyone familiar with this area.',
    'This is more involved than originally estimated — flagging for re-pointing at next planning.',
    'Found a related ticket that overlaps here. Linking them to avoid duplicate work.',
    'Had a quick sync with the design team. They are aligned on the proposed approach.',
    'Dependency resolved — unblocked now and proceeding.'
  ];

  v_comments_b   text[] := ARRAY[
    'Made solid progress. Core logic is in, writing tests now.',
    'PR raised — ready for review. Tagging the relevant folks.',
    'Ran into a gotcha with the existing data model. Added a migration to handle it.',
    'Tests passing locally. Pushing to staging for a smoke test before merging.',
    'Good feedback from the PR review — addressed all comments and pushed updates.',
    'Merged to main. Monitoring for any regressions over the next 24 hours.',
    'Completed and verified on staging. Closing once QA signs off.',
    'Done. Documented the approach in the PR description for future reference.',
    'Wrapped up the first pass. Will need a follow-up ticket for edge cases found along the way.',
    'This turned out cleaner than expected. No breaking changes — safe to ship.'
  ];

  v_comments_c   text[] := ARRAY[
    'One open question: should this be behind a feature flag for the initial rollout?',
    'Note for anyone picking this up: the related config is in the infra repo, not here.',
    'Closing loop on the async discussion — we agreed to go with option B.',
    'Updated the ticket with a revised approach after the team review.',
    'Minor regression found and fixed. Nothing user-facing.',
    'Performance looks good in the profiler — no concerns on that front.',
    'Leaving a note for the next person: check the rate-limit config before deploying to prod.',
    'All edge cases covered. Added regression tests to prevent future breakage.',
    'Confirmed with the product team — scope is locked, no further changes expected.',
    'Final review pending. Should be ready to close by end of sprint.'
  ];

BEGIN
  SELECT ARRAY(
    SELECT id FROM dev_tasks WHERE deleted_at IS NULL ORDER BY id
  ) INTO v_task_ids;

  v_ntasks := coalesce(array_length(v_task_ids, 1), 0);
  IF v_ntasks = 0 THEN
    RAISE EXCEPTION 'No dev tasks found.';
  END IF;

  SELECT ARRAY(SELECT name FROM team_members WHERE deleted_at IS NULL ORDER BY id)
    INTO v_member_names;
  v_nmem := coalesce(array_length(v_member_names, 1), 0);
  IF v_nmem = 0 THEN
    v_member_names := ARRAY['Unassigned'];
    v_nmem := 1;
  END IF;

  RAISE NOTICE 'Adding comments to % dev task(s) with % member(s)…', v_ntasks, v_nmem;

  FOR i IN 1..v_ntasks LOOP
    v_tid := v_task_ids[i];
    v_m1  := v_member_names[1 + floor(random()*v_nmem)::int];
    v_m2  := v_member_names[1 + floor(random()*v_nmem)::int];
    v_m3  := v_member_names[1 + floor(random()*v_nmem)::int];

    INSERT INTO dev_task_comments (created_at, updated_at, task_id, author, content) VALUES
    (NOW() - interval '6 days', NOW() - interval '6 days',
     v_tid, v_m1, v_comments_a[1 + ((i-1) % array_length(v_comments_a,1))]);

    INSERT INTO dev_task_comments (created_at, updated_at, task_id, author, content) VALUES
    (NOW() - interval '3 days', NOW() - interval '3 days',
     v_tid, v_m2, v_comments_b[1 + ((i-1) % array_length(v_comments_b,1))]);

    IF i % 2 = 0 THEN
      INSERT INTO dev_task_comments (created_at, updated_at, task_id, author, content) VALUES
      (NOW() - interval '1 day', NOW() - interval '1 day',
       v_tid, v_m3, v_comments_c[1 + ((i-1) % array_length(v_comments_c,1))]);
    END IF;
  END LOOP;

  RAISE NOTICE 'Done — comments inserted for % dev task(s).', v_ntasks;
END $$;
