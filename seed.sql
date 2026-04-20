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
-- Assignees stored via FK in dev_task_assignees junction table.
-- ─────────────────────────────────────────────────────────────────────────────

DO $$
DECLARE
  v_project_ids  integer[];
  v_member_ids   integer[];
  v_nproj        integer;
  v_nmem         integer;
  v_pid          integer;
  v_m1           integer;
  v_m2           integer;
  v_dt_id        integer;

BEGIN
  SELECT ARRAY(SELECT id FROM projects WHERE deleted_at IS NULL ORDER BY id)
    INTO v_project_ids;
  v_nproj := array_length(v_project_ids, 1);
  IF v_nproj IS NULL THEN RAISE EXCEPTION 'No projects found.'; END IF;

  SELECT ARRAY(SELECT id FROM team_members WHERE deleted_at IS NULL ORDER BY id)
    INTO v_member_ids;
  v_nmem := coalesce(array_length(v_member_ids, 1), 0);
  IF v_nmem = 0 THEN RAISE EXCEPTION 'No team members found — add team members before seeding dev tasks.'; END IF;

  RAISE NOTICE 'Seeding dev tasks across % project(s) with % member(s)…', v_nproj, v_nmem;

  -- ── Improvement tasks ─────────────────────────────────────────────────────
  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_m2 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'38 days', NOW()-interval'10 days', v_pid, 'Migrate REST endpoints to GraphQL', 'Improvement', 'In Progress', 'High')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1), (v_dt_id, v_m2) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'35 days', NOW()-interval'35 days', v_pid, 'Upgrade dashboard charts to Chart.js v4', 'Improvement', 'Todo', 'Medium')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_m2 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'33 days', NOW()-interval'2 days', v_pid, 'Add dark mode support across all pages', 'Improvement', 'In Progress', 'Medium')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1), (v_dt_id, v_m2) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'30 days', NOW()-interval'30 days', v_pid, 'Implement keyboard shortcut navigation', 'Improvement', 'Todo', 'Low')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_m2 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'28 days', NOW()-interval'5 days', v_pid, 'Optimise sprint board rendering performance', 'Improvement', 'Done', 'High')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1), (v_dt_id, v_m2) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'25 days', NOW()-interval'25 days', v_pid, 'Add CSV export for deadline reports', 'Improvement', 'Todo', 'Medium')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_m2 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'22 days', NOW()-interval'8 days', v_pid, 'Redesign onboarding flow for new users', 'Improvement', 'In Progress', 'High')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1), (v_dt_id, v_m2) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'20 days', NOW()-interval'20 days', v_pid, 'Add real-time notification badge via WebSocket', 'Improvement', 'Todo', 'High')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_m2 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'18 days', NOW()-interval'3 days', v_pid, 'Improve pagination UX with infinite scroll option', 'Improvement', 'Done', 'Medium')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1), (v_dt_id, v_m2) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'15 days', NOW()-interval'15 days', v_pid, 'Add multi-language (i18n) support skeleton', 'Improvement', 'Todo', 'Low')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1) ON CONFLICT DO NOTHING;

  -- ── Tech Debt tasks ───────────────────────────────────────────────────────
  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'40 days', NOW()-interval'40 days', v_pid, 'Remove deprecated jQuery dependency from legacy pages', 'Tech Debt', 'Todo', 'High')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_m2 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'37 days', NOW()-interval'12 days', v_pid, 'Consolidate duplicate date-formatting helpers across templates', 'Tech Debt', 'In Progress', 'Medium')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1), (v_dt_id, v_m2) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'34 days', NOW()-interval'4 days', v_pid, 'Replace raw SQL queries in reporting module with GORM', 'Tech Debt', 'Done', 'Critical')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'31 days', NOW()-interval'31 days', v_pid, 'Fix N+1 query on team members listing page', 'Tech Debt', 'Todo', 'High')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_m2 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'27 days', NOW()-interval'9 days', v_pid, 'Add missing database indexes on foreign key columns', 'Tech Debt', 'Done', 'Critical')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1), (v_dt_id, v_m2) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'24 days', NOW()-interval'24 days', v_pid, 'Migrate hardcoded config values to environment variables', 'Tech Debt', 'Todo', 'Medium')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_m2 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'21 days', NOW()-interval'6 days', v_pid, 'Standardise HTTP error response format across all handlers', 'Tech Debt', 'In Progress', 'Medium')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1), (v_dt_id, v_m2) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'17 days', NOW()-interval'17 days', v_pid, 'Remove unused feature-flag code from Sprint 9', 'Tech Debt', 'Todo', 'Low')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_m2 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'14 days', NOW()-interval'14 days', v_pid, 'Increase unit test coverage for auth service to 80%', 'Tech Debt', 'Todo', 'High')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1), (v_dt_id, v_m2) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'10 days', NOW()-interval'2 days', v_pid, 'Upgrade Go module dependencies to latest stable versions', 'Tech Debt', 'In Progress', 'Medium')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1) ON CONFLICT DO NOTHING;

  -- ── Research tasks ────────────────────────────────────────────────────────
  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'36 days', NOW()-interval'36 days', v_pid, 'Evaluate OpenTelemetry for distributed tracing', 'Research', 'Todo', 'Medium')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_m2 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'32 days', NOW()-interval'7 days', v_pid, 'Proof of concept: edge caching with Cloudflare Workers', 'Research', 'In Progress', 'High')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1), (v_dt_id, v_m2) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'29 days', NOW()-interval'29 days', v_pid, 'Assess feasibility of WASM for client-side report rendering', 'Research', 'Todo', 'Low')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_m2 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'26 days', NOW()-interval'11 days', v_pid, 'Benchmark PostgreSQL JSONB vs relational for audit log', 'Research', 'Done', 'Medium')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1), (v_dt_id, v_m2) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'23 days', NOW()-interval'23 days', v_pid, 'Investigate AI-assisted sprint velocity prediction', 'Research', 'Todo', 'Low')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_m2 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'19 days', NOW()-interval'5 days', v_pid, 'Research headless CMS options for documentation portal', 'Research', 'In Progress', 'Medium')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1), (v_dt_id, v_m2) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'16 days', NOW()-interval'16 days', v_pid, 'Explore event-driven architecture with NATS for notifications', 'Research', 'Todo', 'High')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_m2 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'12 days', NOW()-interval'12 days', v_pid, 'Compare ClickHouse vs TimescaleDB for analytics pipeline', 'Research', 'Todo', 'Medium')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1), (v_dt_id, v_m2) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'8 days', NOW()-interval'1 day', v_pid, 'Evaluate Playwright vs Cypress for end-to-end test suite', 'Research', 'In Progress', 'Medium')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1) ON CONFLICT DO NOTHING;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_m2 := v_member_ids[1+floor(random()*v_nmem)::int]; v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO dev_tasks (created_at, updated_at, project_id, title, type, status, priority)
  VALUES (NOW()-interval'5 days', NOW()-interval'5 days', v_pid, 'Spike: feature flags service (LaunchDarkly vs homegrown)', 'Research', 'Todo', 'Low')
  RETURNING id INTO v_dt_id;
  INSERT INTO dev_task_assignees (dev_task_id, team_member_id) VALUES (v_dt_id, v_m1), (v_dt_id, v_m2) ON CONFLICT DO NOTHING;

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

-- ─────────────────────────────────────────────────────────────────────────────
-- Seed: Meeting Minutes
-- ─────────────────────────────────────────────────────────────────────────────

DO $$
DECLARE
  v_project_ids  integer[];
  v_member_names text[];
  v_nproj        integer;
  v_nmem         integer;
  v_pid          integer;
  v_att          text;

BEGIN
  SELECT ARRAY(SELECT id FROM projects WHERE deleted_at IS NULL ORDER BY id) INTO v_project_ids;
  v_nproj := array_length(v_project_ids, 1);
  IF v_nproj IS NULL THEN RAISE EXCEPTION 'No projects found.'; END IF;

  SELECT ARRAY(SELECT name FROM team_members WHERE deleted_at IS NULL ORDER BY id) INTO v_member_names;
  v_nmem := coalesce(array_length(v_member_names, 1), 0);
  IF v_nmem = 0 THEN v_member_names := ARRAY['Unassigned']; v_nmem := 1; END IF;

  RAISE NOTICE 'Seeding meeting minutes across % project(s)…', v_nproj;

  -- Helper: build an attendee CSV from the first few members
  v_att := array_to_string(v_member_names[1:least(4,v_nmem)], ',');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO meetings (created_at, updated_at, project_id, title, date, attendees, notes) VALUES
  (NOW()-interval'38 days', NOW()-interval'38 days', v_pid,
   'Sprint 12 Kick-off Planning',
   'Apr 20, 2026', v_att,
   'Reviewed Sprint 11 retrospective outcomes. Agreed on 54 story points for Sprint 12 focussed on auth module and dashboard scaffold. Story assignment completed. Definition of Done re-confirmed with the team.');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO meetings (created_at, updated_at, project_id, title, date, attendees, notes) VALUES
  (NOW()-interval'36 days', NOW()-interval'36 days', v_pid,
   'Auth Module Design Review',
   'Apr 22, 2026', v_att,
   'Reviewed JWT refresh-token approach vs session cookies. Decided on httpOnly cookie for security. Discussed role-based middleware scope. Agreed to use Redis for session store on staging. Action: Karim to provision Redis by EOD.');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO meetings (created_at, updated_at, project_id, title, date, attendees, notes) VALUES
  (NOW()-interval'34 days', NOW()-interval'34 days', v_pid,
   'Dashboard API Contract Sign-off',
   'Apr 27, 2026', v_att,
   'Product walked through widget requirements. Agreed on 7 endpoints. Pagination strategy confirmed (cursor-based for activity feed, offset for others). Breaking changes policy discussed — versioning via URL prefix.');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO meetings (created_at, updated_at, project_id, title, date, attendees, notes) VALUES
  (NOW()-interval'31 days', NOW()-interval'31 days', v_pid,
   'Notification Service Architecture',
   'Apr 30, 2026', v_att,
   'Decided on polling (30s interval) for v1 notification badge. WebSocket deferred to Sprint 13. Discussed fanout approach for multi-project users. Action: Arif to prototype polling endpoint and benchmark latency.');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO meetings (created_at, updated_at, project_id, title, date, attendees, notes) VALUES
  (NOW()-interval'29 days', NOW()-interval'29 days', v_pid,
   'Sprint 12 Mid-Sprint Check-in',
   'May 01, 2026', v_att,
   'Sprint at 67% — on track. Dashboard widgets all rendering on staging. Chart.js timezone bug mitigated with UTC workaround. Risk flagged: notification badge count may need WebSocket for real-time feel. Decision deferred to end of sprint.');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO meetings (created_at, updated_at, project_id, title, date, attendees, notes) VALUES
  (NOW()-interval'26 days', NOW()-interval'26 days', v_pid,
   'Sprint 12 Retrospective',
   'May 04, 2026', v_att,
   'Went well: async PR reviews, clear API contracts, fast Redis provisioning. Improve: standup timebox felt long — agreed to cap at 10 minutes from Sprint 13. Action: update standup format doc and communicate to team.');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO meetings (created_at, updated_at, project_id, title, date, attendees, notes) VALUES
  (NOW()-interval'24 days', NOW()-interval'24 days', v_pid,
   'Sprint 13 Planning',
   'May 08, 2026', v_att,
   '54 points committed across 3 epics: release pipeline, CI/CD improvements, user settings page. Dependency mapping done. Arif owns CI/CD, Nadia owns settings schema, Rina owns frontend pages. Sprint goal: deploy pipeline live on staging by May 15.');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO meetings (created_at, updated_at, project_id, title, date, attendees, notes) VALUES
  (NOW()-interval'21 days', NOW()-interval'21 days', v_pid,
   'CI/CD Pipeline Design',
   'May 11, 2026', v_att,
   'Reviewed GitHub Actions vs CircleCI — chose GitHub Actions for tighter repo integration. Agreed on two pipelines: PR check (lint + test) and merge-to-main (build + staging deploy). Manual approval gate required for prod. IAM role request filed with DevOps.');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO meetings (created_at, updated_at, project_id, title, date, attendees, notes) VALUES
  (NOW()-interval'18 days', NOW()-interval'18 days', v_pid,
   'User Settings Page UX Review',
   'May 13, 2026', v_att,
   'Walked through Figma mockups for profile, password change, notification prefs, and timezone. Minor feedback: timezone selector should show offset preview. Agreed not to build avatar crop in this sprint. Rina to update mockups and share final by EOD.');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO meetings (created_at, updated_at, project_id, title, date, attendees, notes) VALUES
  (NOW()-interval'15 days', NOW()-interval'15 days', v_pid,
   'Sprint 13 Mid-Sprint Check-in',
   'May 14, 2026', v_att,
   'Sprint at 52%. Timezone refactor flagged as larger than estimated — may slip to Sprint 14. Docker layer caching done, deploy time cut from 4 min to 90s. Settings page on track. No blockers.');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO meetings (created_at, updated_at, project_id, title, date, attendees, notes) VALUES
  (NOW()-interval'12 days', NOW()-interval'12 days', v_pid,
   'v1.4 Release Go/No-Go Meeting',
   'May 21, 2026', v_att,
   'All checks green: regression suite passing (312 tests), staging verified, rollback plan documented, on-call rotation set. One known flaky test in notification suite — agreed to ship and fix post-release. Decision: GO for release on May 22 at 10 AM.');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO meetings (created_at, updated_at, project_id, title, date, attendees, notes) VALUES
  (NOW()-interval'10 days', NOW()-interval'10 days', v_pid,
   'v1.4.1 Hotfix Post-Mortem',
   'May 23, 2026', v_att,
   'Root cause: notification badge not initialising on first login due to missing seeding call in the login handler. Fix was trivial (2 lines). Timeline: reported 11 AM, root cause by 1 PM, deployed 4 PM — 5h total. Action: add login-path integration test to prevent regression.');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO meetings (created_at, updated_at, project_id, title, date, attendees, notes) VALUES
  (NOW()-interval'8 days', NOW()-interval'8 days', v_pid,
   'Sprint 14 Planning',
   'May 25, 2026', v_att,
   '48 points committed: timezone refactor (11 templates), audit log (model + UI), performance profiling. Nadia leads timezone, Arif leads audit log, Rina leads profiling setup. Sprint goal: audit log live and queryable in admin panel by end of sprint.');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO meetings (created_at, updated_at, project_id, title, date, attendees, notes) VALUES
  (NOW()-interval'5 days', NOW()-interval'5 days', v_pid,
   'Audit Log Schema Design',
   'May 26, 2026', v_att,
   'Agreed on schema: id, user_id, action (enum), entity_type, entity_id, metadata (JSONB), created_at. Discussed JSONB vs relational for metadata — chose JSONB for flexibility. Index on (user_id, created_at) confirmed. Arif to write migration and service layer by May 28.');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO meetings (created_at, updated_at, project_id, title, date, attendees, notes) VALUES
  (NOW()-interval'3 days', NOW()-interval'3 days', v_pid,
   'Sprint 14 Mid-Sprint Check-in',
   'May 28, 2026', v_att,
   'Sprint at 61%. Timezone refactor fully merged. Audit log query optimised to 18ms after composite index. Rina profiling 3 more endpoints this week. No blockers. On track for audit log UI by end of sprint.');

  RAISE NOTICE 'Done — 15 meeting entries inserted.';
END $$;

-- ─────────────────────────────────────────────────────────────────────────────
-- Seed: Slack Threads
-- ─────────────────────────────────────────────────────────────────────────────

DO $$
DECLARE
  v_project_ids  integer[];
  v_member_ids   integer[];
  v_nproj        integer;
  v_nmem         integer;
  v_pid          integer;
  v_author_id    integer;

BEGIN
  SELECT ARRAY(SELECT id FROM projects WHERE deleted_at IS NULL ORDER BY id) INTO v_project_ids;
  v_nproj := array_length(v_project_ids, 1);
  IF v_nproj IS NULL THEN RAISE EXCEPTION 'No projects found.'; END IF;

  SELECT ARRAY(SELECT id FROM team_members WHERE deleted_at IS NULL ORDER BY id) INTO v_member_ids;
  v_nmem := coalesce(array_length(v_member_ids, 1), 0);
  IF v_nmem = 0 THEN RAISE EXCEPTION 'No team members found.'; END IF;

  RAISE NOTICE 'Seeding slack threads…';

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int]; v_author_id := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO slack_threads (created_at, updated_at, project_id, message_link, topic, summary, tags, author_id) VALUES
  (NOW()-interval'35 days', NOW()-interval'35 days', v_pid,
   'https://slack.com/archives/C01/p1714000001',
   'JWT refresh token expiry strategy',
   'Discussed whether to use sliding window or fixed expiry for refresh tokens. Team aligned on 7-day fixed expiry with silent renewal on activity. Rolling tokens ruled out due to complexity with multi-device sessions.',
   'auth,security,backend', v_author_id);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int]; v_author_id := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO slack_threads (created_at, updated_at, project_id, message_link, topic, summary, tags, author_id) VALUES
  (NOW()-interval'33 days', NOW()-interval'33 days', v_pid,
   'https://slack.com/archives/C01/p1714000002',
   'Redis session store — staging provisioning',
   'Karim confirmed Redis 7.2 provisioned on staging. Connection string shared in #infra-secrets. Session invalidation tests unblocked. No persistence needed for staging; AOF enabled for prod.',
   'infra,redis,staging', v_author_id);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int]; v_author_id := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO slack_threads (created_at, updated_at, project_id, message_link, topic, summary, tags, author_id) VALUES
  (NOW()-interval'31 days', NOW()-interval'31 days', v_pid,
   'https://slack.com/archives/C01/p1714000003',
   'Dashboard widget API contract — product feedback',
   'Product approved all 7 widget endpoints. Requested one change: sprint velocity chart should include carry-over points separately. Nadia to update the response schema. No timeline impact.',
   'api,dashboard,product', v_author_id);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int]; v_author_id := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO slack_threads (created_at, updated_at, project_id, message_link, topic, summary, tags, author_id) VALUES
  (NOW()-interval'29 days', NOW()-interval'29 days', v_pid,
   'https://slack.com/archives/C01/p1714000004',
   'Chart.js v3 timezone offset bug',
   'Confirmed upstream Chart.js issue #12350. Workaround: force dataset labels to UTC strings and display local time in tooltips only. Fix expected in Chart.js v4.1 — upgrade tracked in dev tasks.',
   'frontend,bug,chartsjs', v_author_id);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int]; v_author_id := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO slack_threads (created_at, updated_at, project_id, message_link, topic, summary, tags, author_id) VALUES
  (NOW()-interval'27 days', NOW()-interval'27 days', v_pid,
   'https://slack.com/archives/C01/p1714000005',
   'N+1 query fix on sprint-summary endpoint',
   'Root cause: tasks loaded one-by-one inside a loop. Fix: added Preload("Tasks") to sprint repo. Query time dropped from 340ms to 28ms on staging dataset (800 tasks). Added benchmark test to prevent regression.',
   'performance,backend,gorm', v_author_id);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int]; v_author_id := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO slack_threads (created_at, updated_at, project_id, message_link, topic, summary, tags, author_id) VALUES
  (NOW()-interval'25 days', NOW()-interval'25 days', v_pid,
   'https://slack.com/archives/C01/p1714000006',
   'Notification polling vs WebSocket trade-off',
   'Agreed on 30s polling for Sprint 12 given timeline constraints. WebSocket (socket.io or Go gorilla/websocket) deferred to Sprint 13. Polling endpoint should use ETag caching to reduce payload on no-change responses.',
   'notifications,architecture,websocket', v_author_id);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int]; v_author_id := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO slack_threads (created_at, updated_at, project_id, message_link, topic, summary, tags, author_id) VALUES
  (NOW()-interval'22 days', NOW()-interval'22 days', v_pid,
   'https://slack.com/archives/C01/p1714000007',
   'GitHub Actions IAM role for staging deploys',
   'New IAM role arn:aws:iam::123456789:role/github-actions-staging approved by Karim. Least-privilege: S3 write to deploy bucket + ECS task update only. Secrets stored in GitHub repo environment "staging".',
   'ci-cd,infra,devops', v_author_id);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int]; v_author_id := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO slack_threads (created_at, updated_at, project_id, message_link, topic, summary, tags, author_id) VALUES
  (NOW()-interval'20 days', NOW()-interval'20 days', v_pid,
   'https://slack.com/archives/C01/p1714000008',
   'Docker layer caching strategy for faster builds',
   'Arif restructured Dockerfile: dependencies layer first, source copy last. Added --mount=type=cache for Go module cache. Staging deploy time: 4min → 1min 40s. Template shared in #engineering for other services.',
   'ci-cd,docker,performance', v_author_id);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int]; v_author_id := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO slack_threads (created_at, updated_at, project_id, message_link, topic, summary, tags, author_id) VALUES
  (NOW()-interval'17 days', NOW()-interval'17 days', v_pid,
   'https://slack.com/archives/C01/p1714000009',
   'Avatar upload file size validation',
   'Bug: silent failure on files >2MB. Root cause: missing server-side validation (relied on frontend only). Fix: added MaxBytesReader in handler + user-facing error message. Also added accepted MIME type check (image/jpeg, image/png, image/webp).',
   'bug,backend,ux', v_author_id);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int]; v_author_id := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO slack_threads (created_at, updated_at, project_id, message_link, topic, summary, tags, author_id) VALUES
  (NOW()-interval'15 days', NOW()-interval'15 days', v_pid,
   'https://slack.com/archives/C01/p1714000010',
   'Prod deployment approval gate design',
   'Manual approval gate implemented via GitHub Actions environment protection rules. Required reviewers: any 1 of [lead, devops]. Timeout: 2 hours before auto-cancel. Prod secrets isolated in "production" environment.',
   'ci-cd,devops,release', v_author_id);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int]; v_author_id := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO slack_threads (created_at, updated_at, project_id, message_link, topic, summary, tags, author_id) VALUES
  (NOW()-interval'13 days', NOW()-interval'13 days', v_pid,
   'https://slack.com/archives/C01/p1714000011',
   'v1.4 release comms and rollback plan',
   'Release notes drafted and shared in #product-updates. Rollback plan: ECS task definition rollback to previous revision (< 5 min). DB migrations are additive-only so no rollback needed. On-call: Arif primary, Karim secondary.',
   'release,comms,ops', v_author_id);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int]; v_author_id := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO slack_threads (created_at, updated_at, project_id, message_link, topic, summary, tags, author_id) VALUES
  (NOW()-interval'10 days', NOW()-interval'10 days', v_pid,
   'https://slack.com/archives/C01/p1714000012',
   'Notification badge init bug root cause',
   'Bug in login handler: badge count seeded from DB after session cookie set but before response flushed — race condition on first page load. Fix: move badge seed to session middleware. Hotfix deployed v1.4.1 same day.',
   'bug,notifications,hotfix', v_author_id);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int]; v_author_id := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO slack_threads (created_at, updated_at, project_id, message_link, topic, summary, tags, author_id) VALUES
  (NOW()-interval'8 days', NOW()-interval'8 days', v_pid,
   'https://slack.com/archives/C01/p1714000013',
   'Audit log JSONB vs relational metadata',
   'Benchmarked both approaches on 1M row dataset. JSONB: flexible schema, 12ms read, 8ms write. Relational: rigid schema, 9ms read, 6ms write. Decision: JSONB for metadata given evolving audit requirements. Index on (user_id, created_at) added.',
   'database,architecture,audit', v_author_id);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int]; v_author_id := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO slack_threads (created_at, updated_at, project_id, message_link, topic, summary, tags, author_id) VALUES
  (NOW()-interval'5 days', NOW()-interval'5 days', v_pid,
   'https://slack.com/archives/C01/p1714000014',
   'Timezone display format — UX decision',
   'UX call outcome: show absolute local time (e.g. "2 May, 14:30") with UTC offset in tooltip on hover. Relative time ("2 hours ago") shown only in activity feeds. Rina to update component library. Nadia to start template refactor.',
   'ux,frontend,timezone', v_author_id);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int]; v_author_id := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO slack_threads (created_at, updated_at, project_id, message_link, topic, summary, tags, author_id) VALUES
  (NOW()-interval'3 days', NOW()-interval'3 days', v_pid,
   'https://slack.com/archives/C01/p1714000015',
   'Performance profiling findings — Sprint 14',
   'pprof results: top 3 slow endpoints are audit log list (fixed), sprint board load (340ms, index missing on sprint_tasks.sprint_id), and release detail (220ms, N+1 on stages). Two new tickets created for Sprint 14 backlog.',
   'performance,profiling,backend', v_author_id);

  RAISE NOTICE 'Done — 15 slack thread entries inserted.';
END $$;

-- ─────────────────────────────────────────────────────────────────────────────
-- Seed: Sticky Notes
-- ─────────────────────────────────────────────────────────────────────────────

DO $$
DECLARE
  v_project_ids integer[];
  v_nproj       integer;
  v_pid         integer;

BEGIN
  SELECT ARRAY(SELECT id FROM projects WHERE deleted_at IS NULL ORDER BY id) INTO v_project_ids;
  v_nproj := array_length(v_project_ids, 1);
  IF v_nproj IS NULL THEN RAISE EXCEPTION 'No projects found.'; END IF;

  RAISE NOTICE 'Seeding sticky notes…';

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sticky_notes (created_at, updated_at, project_id, title, content, color, pinned) VALUES
  (NOW()-interval'30 days', NOW()-interval'30 days', v_pid,
   'Sprint 12 Goals',
   'Auth module (JWT + refresh tokens), dashboard scaffold, 4 carry-over tickets from Sprint 11. Target: 54 points. Definition of Done: code reviewed, tests green, deployed to staging.',
   'yellow', true);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sticky_notes (created_at, updated_at, project_id, title, content, color, pinned) VALUES
  (NOW()-interval'28 days', NOW()-interval'28 days', v_pid,
   'Redis config (staging)',
   'Host: redis-staging.internal:6379. No auth (internal only). No persistence. Max memory: 256MB, policy: allkeys-lru. Use DB 0 for sessions, DB 1 for rate limiting.',
   'blue', false);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sticky_notes (created_at, updated_at, project_id, title, content, color, pinned) VALUES
  (NOW()-interval'26 days', NOW()-interval'26 days', v_pid,
   'API versioning convention',
   'All breaking changes go under /v2/ prefix. Non-breaking additions are backwards-compatible in /v1/. Deprecation notice required 2 sprints before removal. Document changes in CHANGELOG.md.',
   'green', false);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sticky_notes (created_at, updated_at, project_id, title, content, color, pinned) VALUES
  (NOW()-interval'24 days', NOW()-interval'24 days', v_pid,
   'PR checklist',
   '[ ] Tests written and passing. [ ] No N+1 queries introduced. [ ] Error messages are user-friendly. [ ] Migrations are reversible. [ ] Secrets not committed. [ ] PR description explains the why.',
   'pink', true);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sticky_notes (created_at, updated_at, project_id, title, content, color, pinned) VALUES
  (NOW()-interval'22 days', NOW()-interval'22 days', v_pid,
   'On-call rotation',
   'Week 1: Arif. Week 2: Nadia. Week 3: Karim. Week 4: Rina. Escalation: Slack #incidents first, then PagerDuty if no ack in 10 min. Runbook: notion.so/team/runbooks.',
   'yellow', false);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sticky_notes (created_at, updated_at, project_id, title, content, color, pinned) VALUES
  (NOW()-interval'20 days', NOW()-interval'20 days', v_pid,
   'Notification polling endpoint',
   'GET /api/notifications/count — returns unread count. Use ETag for caching. Client polls every 30s. Upgrade to WebSocket planned for Sprint 13. Rate limit: 10 req/min per user.',
   'blue', false);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sticky_notes (created_at, updated_at, project_id, title, content, color, pinned) VALUES
  (NOW()-interval'18 days', NOW()-interval'18 days', v_pid,
   'Sprint 13 Goals',
   'Release pipeline (staging + prod), CI/CD with GitHub Actions, user settings page (profile, password, notification prefs, timezone). 54 points. Sprint goal: prod pipeline live by May 15.',
   'yellow', true);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sticky_notes (created_at, updated_at, project_id, title, content, color, pinned) VALUES
  (NOW()-interval'16 days', NOW()-interval'16 days', v_pid,
   'Staging deploy checklist',
   '1. Merge to main. 2. GitHub Actions triggers automatically. 3. Wait for build (~90s). 4. Check ECS health — green? 5. Smoke test: /health endpoint + login flow. 6. Post in #deployments.',
   'green', false);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sticky_notes (created_at, updated_at, project_id, title, content, color, pinned) VALUES
  (NOW()-interval'14 days', NOW()-interval'14 days', v_pid,
   'Accessibility findings (Sprint 13)',
   '3 issues found in settings page audit: 1) Missing aria-label on icon buttons. 2) Focus trap broken in modal on Safari. 3) Color contrast ratio 3.2:1 on placeholder text (needs 4.5:1). All fixed before release.',
   'pink', false);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sticky_notes (created_at, updated_at, project_id, title, content, color, pinned) VALUES
  (NOW()-interval'12 days', NOW()-interval'12 days', v_pid,
   'v1.4 release checklist',
   '[ ] Regression suite green (312 tests). [ ] Staging verified by QA. [ ] Release notes approved. [ ] Rollback plan documented. [ ] On-call set. [ ] DB migrations additive only. [ ] Monitoring dashboards reviewed.',
   'yellow', false);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sticky_notes (created_at, updated_at, project_id, title, content, color, pinned) VALUES
  (NOW()-interval'10 days', NOW()-interval'10 days', v_pid,
   'Prod DB connection pool',
   'Increased from 10 → 25 connections before v1.4 release. Monitor pg_stat_activity during peak. Max connections on RDS instance: 100. Leave 25% headroom for migrations and admin connections.',
   'blue', false);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sticky_notes (created_at, updated_at, project_id, title, content, color, pinned) VALUES
  (NOW()-interval'8 days', NOW()-interval'8 days', v_pid,
   'Sprint 14 Goals',
   'Timezone refactor (11 templates), audit log (model + service + UI), performance profiling. 48 points. Key risk: timezone refactor touches many files — run full regression before merging.',
   'yellow', true);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sticky_notes (created_at, updated_at, project_id, title, content, color, pinned) VALUES
  (NOW()-interval'6 days', NOW()-interval'6 days', v_pid,
   'Audit log schema',
   'Table: audit_logs. Columns: id, user_id (FK), action (enum: create/update/delete/login/logout), entity_type (varchar), entity_id (uint), metadata (JSONB), created_at. Index: (user_id, created_at). No soft delete.',
   'green', false);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sticky_notes (created_at, updated_at, project_id, title, content, color, pinned) VALUES
  (NOW()-interval'4 days', NOW()-interval'4 days', v_pid,
   'Performance profiling notes',
   'Slow endpoints: sprint_board (340ms) — missing index on sprint_tasks.sprint_id. release_detail (220ms) — N+1 on stages. Both ticketed for Sprint 14. Use pprof at /debug/pprof/ (admin only).',
   'pink', false);

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sticky_notes (created_at, updated_at, project_id, title, content, color, pinned) VALUES
  (NOW()-interval'2 days', NOW()-interval'2 days', v_pid,
   'Team Slack channels',
   '#engineering — general dev discussion. #incidents — alerts and on-call. #deployments — deploy logs. #product-updates — release comms. #infra-secrets — credentials (restricted). #standup — async updates.',
   'purple', false);

  RAISE NOTICE 'Done — 15 sticky note entries inserted.';
END $$;

-- ─────────────────────────────────────────────────────────────────────────────
-- Seed: Deadlines
-- ─────────────────────────────────────────────────────────────────────────────

DO $$
DECLARE
  v_project_ids integer[];
  v_nproj       integer;
  v_pid         integer;

BEGIN
  SELECT ARRAY(SELECT id FROM projects WHERE deleted_at IS NULL ORDER BY id) INTO v_project_ids;
  v_nproj := array_length(v_project_ids, 1);
  IF v_nproj IS NULL THEN RAISE EXCEPTION 'No projects found.'; END IF;

  RAISE NOTICE 'Seeding deadlines…';

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO deadlines (created_at, updated_at, project_id, title, due_date, priority) VALUES
  (NOW()-interval'30 days', NOW()-interval'30 days', v_pid, 'Auth module code freeze', '2026-05-02', 'Critical');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO deadlines (created_at, updated_at, project_id, title, due_date, priority) VALUES
  (NOW()-interval'28 days', NOW()-interval'28 days', v_pid, 'Dashboard API contract sign-off', '2026-05-05', 'High');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO deadlines (created_at, updated_at, project_id, title, due_date, priority) VALUES
  (NOW()-interval'25 days', NOW()-interval'25 days', v_pid, 'Sprint 12 QA pass', '2026-05-07', 'High');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO deadlines (created_at, updated_at, project_id, title, due_date, priority) VALUES
  (NOW()-interval'22 days', NOW()-interval'22 days', v_pid, 'v1.4 release candidate tag', '2026-05-20', 'Critical');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO deadlines (created_at, updated_at, project_id, title, due_date, priority) VALUES
  (NOW()-interval'20 days', NOW()-interval'20 days', v_pid, 'v1.4 production release', '2026-05-22', 'Critical');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO deadlines (created_at, updated_at, project_id, title, due_date, priority) VALUES
  (NOW()-interval'18 days', NOW()-interval'18 days', v_pid, 'CI/CD staging pipeline live', '2026-05-15', 'High');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO deadlines (created_at, updated_at, project_id, title, due_date, priority) VALUES
  (NOW()-interval'15 days', NOW()-interval'15 days', v_pid, 'User settings page accessibility audit', '2026-05-18', 'Medium');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO deadlines (created_at, updated_at, project_id, title, due_date, priority) VALUES
  (NOW()-interval'12 days', NOW()-interval'12 days', v_pid, 'Timezone refactor complete', '2026-05-27', 'High');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO deadlines (created_at, updated_at, project_id, title, due_date, priority) VALUES
  (NOW()-interval'10 days', NOW()-interval'10 days', v_pid, 'Audit log UI in admin panel', '2026-06-06', 'High');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO deadlines (created_at, updated_at, project_id, title, due_date, priority) VALUES
  (NOW()-interval'8 days', NOW()-interval'8 days', v_pid, 'Performance profiling report', '2026-06-01', 'Medium');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO deadlines (created_at, updated_at, project_id, title, due_date, priority) VALUES
  (NOW()-interval'6 days', NOW()-interval'6 days', v_pid, 'Sprint 14 retrospective', '2026-06-12', 'Low');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO deadlines (created_at, updated_at, project_id, title, due_date, priority) VALUES
  (NOW()-interval'4 days', NOW()-interval'4 days', v_pid, 'Legal compliance review — session storage', '2026-06-20', 'Critical');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO deadlines (created_at, updated_at, project_id, title, due_date, priority) VALUES
  (NOW()-interval'3 days', NOW()-interval'3 days', v_pid, 'WebSocket spike completion', '2026-06-15', 'Medium');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO deadlines (created_at, updated_at, project_id, title, due_date, priority) VALUES
  (NOW()-interval'2 days', NOW()-interval'2 days', v_pid, 'v1.5 scope finalisation', '2026-06-25', 'High');

  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO deadlines (created_at, updated_at, project_id, title, due_date, priority) VALUES
  (NOW()-interval'1 day', NOW()-interval'1 day', v_pid, 'Sprint 15 planning session', '2026-06-30', 'Medium');

  RAISE NOTICE 'Done — 15 deadline entries inserted.';
END $$;

-- ─────────────────────────────────────────────────────────────────────────────
-- Seed: Sprints (with tasks)
-- Assignees stored via FK in sprint_task_assignees junction table.
-- ─────────────────────────────────────────────────────────────────────────────

DO $$
DECLARE
  v_project_ids  integer[];
  v_member_ids   integer[];
  v_nproj        integer;
  v_nmem         integer;
  v_pid          integer;
  v_m1           integer;
  v_m2           integer;
  v_sprint_id    integer;
  v_task_id      integer;

BEGIN
  SELECT ARRAY(SELECT id FROM projects WHERE deleted_at IS NULL ORDER BY id) INTO v_project_ids;
  v_nproj := array_length(v_project_ids, 1);
  IF v_nproj IS NULL THEN RAISE EXCEPTION 'No projects found.'; END IF;

  SELECT ARRAY(SELECT id FROM team_members WHERE deleted_at IS NULL ORDER BY id) INTO v_member_ids;
  v_nmem := coalesce(array_length(v_member_ids, 1), 0);
  IF v_nmem = 0 THEN RAISE EXCEPTION 'No team members found — add team members before seeding sprints.'; END IF;

  RAISE NOTICE 'Seeding sprints and tasks…';

  -- ── Sprint 11 (completed) ───────────────────────────────────────────────────
  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sprints (created_at, updated_at, project_id, name, goal, progress, start_date, end_date, active)
  VALUES (NOW()-interval'50 days', NOW()-interval'36 days', v_pid,
    'Sprint 11', 'Complete user profile API and fix reported bugs from v1.3', 100,
    '2026-04-06', '2026-04-18', false)
  RETURNING id INTO v_sprint_id;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_m2 := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'50 days', NOW()-interval'38 days', v_sprint_id, 'User profile GET/PUT endpoints', 'Done', 'High')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m1) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'49 days', NOW()-interval'39 days', v_sprint_id, 'Fix pagination off-by-one on activity feed', 'Done', 'Medium')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'48 days', NOW()-interval'40 days', v_sprint_id, 'Avatar upload endpoint', 'Done', 'Medium')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m1) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'47 days', NOW()-interval'41 days', v_sprint_id, 'Email uniqueness validation on profile update', 'Done', 'Low')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'46 days', NOW()-interval'42 days', v_sprint_id, 'Fix XSS in user display name rendering', 'Done', 'Critical')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m1), (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'45 days', NOW()-interval'43 days', v_sprint_id, 'Write service-layer tests for profile module', 'Done', 'Medium')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  -- ── Sprint 12 (completed) ───────────────────────────────────────────────────
  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sprints (created_at, updated_at, project_id, name, goal, progress, start_date, end_date, active)
  VALUES (NOW()-interval'35 days', NOW()-interval'22 days', v_pid,
    'Sprint 12', 'Ship auth module, dashboard scaffold, and in-app notifications', 100,
    '2026-04-20', '2026-05-08', false)
  RETURNING id INTO v_sprint_id;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_m2 := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'35 days', NOW()-interval'30 days', v_sprint_id, 'JWT refresh-token endpoint', 'Done', 'High')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m1) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'34 days', NOW()-interval'29 days', v_sprint_id, 'Role-based access middleware', 'Done', 'High')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'33 days', NOW()-interval'28 days', v_sprint_id, 'Password reset email flow', 'Done', 'Medium')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m1) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'32 days', NOW()-interval'27 days', v_sprint_id, 'Dashboard top nav and sidebar', 'Done', 'Medium')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'31 days', NOW()-interval'26 days', v_sprint_id, 'Dashboard widget API endpoints (7)', 'Done', 'High')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m1), (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'30 days', NOW()-interval'25 days', v_sprint_id, 'In-app notification model and dispatch', 'Done', 'High')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m1) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'29 days', NOW()-interval'24 days', v_sprint_id, 'Notification badge polling (30s)', 'Done', 'Medium')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'28 days', NOW()-interval'23 days', v_sprint_id, 'Fix Chart.js timezone offset bug', 'Done', 'Medium')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'27 days', NOW()-interval'22 days', v_sprint_id, 'Auth integration tests', 'Done', 'High')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m1) ON CONFLICT DO NOTHING;

  -- ── Sprint 13 (completed) ───────────────────────────────────────────────────
  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sprints (created_at, updated_at, project_id, name, goal, progress, start_date, end_date, active)
  VALUES (NOW()-interval'21 days', NOW()-interval'8 days', v_pid,
    'Sprint 13', 'Release pipeline live on staging and prod; user settings page complete', 100,
    '2026-05-11', '2026-05-22', false)
  RETURNING id INTO v_sprint_id;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_m2 := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'21 days', NOW()-interval'16 days', v_sprint_id, 'GitHub Actions PR check pipeline', 'Done', 'High')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m1) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'20 days', NOW()-interval'15 days', v_sprint_id, 'Staging auto-deploy on merge to main', 'Done', 'High')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m1) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'19 days', NOW()-interval'14 days', v_sprint_id, 'Prod pipeline with manual approval gate', 'Done', 'Critical')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m1) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'18 days', NOW()-interval'13 days', v_sprint_id, 'Docker layer caching (4min → 90s)', 'Done', 'Medium')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m1) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'17 days', NOW()-interval'12 days', v_sprint_id, 'User settings schema and migration', 'Done', 'High')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'16 days', NOW()-interval'11 days', v_sprint_id, 'Profile and password change endpoints', 'Done', 'High')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'15 days', NOW()-interval'10 days', v_sprint_id, 'Notification preferences endpoint', 'Done', 'Medium')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'14 days', NOW()-interval'9 days', v_sprint_id, 'Settings page UI (Figma → implementation)', 'Done', 'High')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'13 days', NOW()-interval'8 days', v_sprint_id, 'Accessibility audit fixes (3 issues)', 'Done', 'Medium')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  -- ── Sprint 14 (active) ──────────────────────────────────────────────────────
  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  INSERT INTO sprints (created_at, updated_at, project_id, name, goal, progress, start_date, end_date, active)
  VALUES (NOW()-interval'7 days', NOW(), v_pid,
    'Sprint 14', 'Timezone refactor complete, audit log live in admin panel, performance profiling report', 61,
    '2026-05-25', '2026-06-06', true)
  RETURNING id INTO v_sprint_id;

  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int]; v_m2 := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'7 days', NOW()-interval'2 days', v_sprint_id, 'Timezone refactor — 11 templates', 'Done', 'High')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m1) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'7 days', NOW()-interval'3 days', v_sprint_id, 'Audit log model and migration', 'Done', 'High')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'6 days', NOW()-interval'2 days', v_sprint_id, 'Composite index on (user_id, created_at)', 'Done', 'High')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'6 days', NOW()-interval'1 day', v_sprint_id, 'Audit log service layer + tests', 'In Progress', 'High')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'5 days', NOW()-interval'1 day', v_sprint_id, 'Audit log admin UI (filterable table)', 'In Progress', 'Medium')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'5 days', NOW()-interval'5 days', v_sprint_id, 'pprof endpoint (admin only)', 'Done', 'Medium')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m1) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'4 days', NOW()-interval'1 day', v_sprint_id, 'Fix N+1 on sprint board load', 'In Progress', 'High')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m1) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'4 days', NOW()-interval'4 days', v_sprint_id, 'Fix N+1 on release detail stages', 'Todo', 'Medium')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m1) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'3 days', NOW()-interval'3 days', v_sprint_id, 'Performance profiling report (3 endpoints)', 'Todo', 'Medium')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m1) ON CONFLICT DO NOTHING;

  INSERT INTO sprint_tasks (created_at, updated_at, sprint_id, title, status, priority)
  VALUES (NOW()-interval'2 days', NOW()-interval'2 days', v_sprint_id, 'Fix flaky notification test (retry logic)', 'Todo', 'Low')
  RETURNING id INTO v_task_id;
  INSERT INTO sprint_task_assignees (sprint_task_id, team_member_id) VALUES (v_task_id, v_m2) ON CONFLICT DO NOTHING;

  RAISE NOTICE 'Done — 4 sprints and tasks inserted.';
END $$;

-- ─────────────────────────────────────────────────────────────────────────────
-- Seed: Releases (with stages and stories)
-- ─────────────────────────────────────────────────────────────────────────────

DO $$
DECLARE
  v_project_ids  integer[];
  v_member_ids   integer[];
  v_nproj        integer;
  v_nmem         integer;
  v_pid          integer;
  v_m1           integer;
  v_m2           integer;
  v_release_id   integer;
  v_stage_id     integer;

BEGIN
  SELECT ARRAY(SELECT id FROM projects WHERE deleted_at IS NULL ORDER BY id) INTO v_project_ids;
  v_nproj := array_length(v_project_ids, 1);
  IF v_nproj IS NULL THEN RAISE EXCEPTION 'No projects found.'; END IF;

  SELECT ARRAY(SELECT id FROM team_members WHERE deleted_at IS NULL ORDER BY id) INTO v_member_ids;
  v_nmem := coalesce(array_length(v_member_ids, 1), 0);
  IF v_nmem = 0 THEN RAISE EXCEPTION 'No team members found — add team members before seeding releases.'; END IF;

  RAISE NOTICE 'Seeding releases…';

  -- ── v1.3 — Released ────────────────────────────────────────────────────────
  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int];
  v_m2 := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO releases (created_at, updated_at, project_id, name, description, status, target_date)
  VALUES (NOW()-interval'55 days', NOW()-interval'40 days', v_pid,
    'v1.3.0', 'Bug fixes from v1.2 feedback: pagination, XSS fix, profile API improvements.', 'Released', '2026-04-18')
  RETURNING id INTO v_release_id;

  INSERT INTO release_stages (created_at, updated_at, release_id, name, status)
  VALUES (NOW()-interval'55 days', NOW()-interval'50 days', v_release_id, 'Development', 'Done')
  RETURNING id INTO v_stage_id;
  INSERT INTO release_stories (created_at, updated_at, stage_id, title, assignee_id, status) VALUES
  (NOW()-interval'54 days', NOW()-interval'50 days', v_stage_id, 'Fix pagination off-by-one on activity feed', v_m1, 'Passed'),
  (NOW()-interval'53 days', NOW()-interval'50 days', v_stage_id, 'Fix XSS in user display name', v_m2, 'Passed'),
  (NOW()-interval'52 days', NOW()-interval'50 days', v_stage_id, 'User profile GET/PUT endpoints', v_m1, 'Passed');

  INSERT INTO release_stages (created_at, updated_at, release_id, name, status)
  VALUES (NOW()-interval'49 days', NOW()-interval'43 days', v_release_id, 'QA', 'Done')
  RETURNING id INTO v_stage_id;
  INSERT INTO release_stories (created_at, updated_at, stage_id, title, assignee_id, status) VALUES
  (NOW()-interval'48 days', NOW()-interval'43 days', v_stage_id, 'Regression test suite run', v_m2, 'Passed'),
  (NOW()-interval'47 days', NOW()-interval'43 days', v_stage_id, 'Smoke test on staging', v_m1, 'Passed');

  INSERT INTO release_stages (created_at, updated_at, release_id, name, status)
  VALUES (NOW()-interval'42 days', NOW()-interval'40 days', v_release_id, 'Production Deploy', 'Done')
  RETURNING id INTO v_stage_id;
  INSERT INTO release_stories (created_at, updated_at, stage_id, title, assignee_id, status) VALUES
  (NOW()-interval'41 days', NOW()-interval'40 days', v_stage_id, 'Deploy to production', v_m1, 'Passed'),
  (NOW()-interval'40 days', NOW()-interval'40 days', v_stage_id, 'Post-deploy health check', v_m2, 'Passed');

  -- ── v1.4 — Released ────────────────────────────────────────────────────────
  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int];
  v_m2 := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO releases (created_at, updated_at, project_id, name, description, status, target_date)
  VALUES (NOW()-interval'35 days', NOW()-interval'11 days', v_pid,
    'v1.4.0', 'Auth module, dashboard with widgets, in-app notifications, user settings page.', 'Released', '2026-05-22')
  RETURNING id INTO v_release_id;

  INSERT INTO release_stages (created_at, updated_at, release_id, name, status)
  VALUES (NOW()-interval'35 days', NOW()-interval'22 days', v_release_id, 'Development', 'Done')
  RETURNING id INTO v_stage_id;
  INSERT INTO release_stories (created_at, updated_at, stage_id, title, assignee_id, status) VALUES
  (NOW()-interval'34 days', NOW()-interval'25 days', v_stage_id, 'JWT refresh-token + role middleware', v_m1, 'Passed'),
  (NOW()-interval'33 days', NOW()-interval'26 days', v_stage_id, 'Dashboard widget API (7 endpoints)', v_m2, 'Passed'),
  (NOW()-interval'32 days', NOW()-interval'27 days', v_stage_id, 'In-app notifications with polling', v_m1, 'Passed'),
  (NOW()-interval'31 days', NOW()-interval'28 days', v_stage_id, 'User settings page (profile + password + prefs)', v_m2, 'Passed'),
  (NOW()-interval'30 days', NOW()-interval'29 days', v_stage_id, 'CI/CD pipeline (staging + prod)', v_m1, 'Passed');

  INSERT INTO release_stages (created_at, updated_at, release_id, name, status)
  VALUES (NOW()-interval'21 days', NOW()-interval'13 days', v_release_id, 'QA', 'Done')
  RETURNING id INTO v_stage_id;
  INSERT INTO release_stories (created_at, updated_at, stage_id, title, assignee_id, status) VALUES
  (NOW()-interval'20 days', NOW()-interval'14 days', v_stage_id, 'Full regression suite (312 tests)', v_m2, 'Passed'),
  (NOW()-interval'19 days', NOW()-interval'14 days', v_stage_id, 'Accessibility audit — settings page', v_m2, 'Passed'),
  (NOW()-interval'18 days', NOW()-interval'13 days', v_stage_id, 'Exploratory testing — dashboard + notifications', v_m1, 'Passed');

  INSERT INTO release_stages (created_at, updated_at, release_id, name, status)
  VALUES (NOW()-interval'12 days', NOW()-interval'11 days', v_release_id, 'Production Deploy', 'Done')
  RETURNING id INTO v_stage_id;
  INSERT INTO release_stories (created_at, updated_at, stage_id, title, assignee_id, status) VALUES
  (NOW()-interval'12 days', NOW()-interval'11 days', v_stage_id, 'Zero-downtime deploy to production', v_m1, 'Passed'),
  (NOW()-interval'11 days', NOW()-interval'11 days', v_stage_id, 'Post-deploy smoke test and monitoring check', v_m2, 'Passed');

  -- ── v1.4.1 — Hotfix Released ────────────────────────────────────────────────
  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int];
  v_m2 := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO releases (created_at, updated_at, project_id, name, description, status, target_date)
  VALUES (NOW()-interval'10 days', NOW()-interval'9 days', v_pid,
    'v1.4.1', 'Hotfix: notification badge not initialising on first login after v1.4.0 release.', 'Released', '2026-05-22')
  RETURNING id INTO v_release_id;

  INSERT INTO release_stages (created_at, updated_at, release_id, name, status)
  VALUES (NOW()-interval'10 days', NOW()-interval'10 days', v_release_id, 'Hotfix Development', 'Done')
  RETURNING id INTO v_stage_id;
  INSERT INTO release_stories (created_at, updated_at, stage_id, title, assignee_id, status) VALUES
  (NOW()-interval'10 days', NOW()-interval'10 days', v_stage_id, 'Move badge seed to session middleware', v_m1, 'Passed'),
  (NOW()-interval'10 days', NOW()-interval'10 days', v_stage_id, 'Add login-path integration test', v_m2, 'Passed');

  INSERT INTO release_stages (created_at, updated_at, release_id, name, status)
  VALUES (NOW()-interval'9 days', NOW()-interval'9 days', v_release_id, 'Production Hotfix Deploy', 'Done')
  RETURNING id INTO v_stage_id;
  INSERT INTO release_stories (created_at, updated_at, stage_id, title, assignee_id, status) VALUES
  (NOW()-interval'9 days', NOW()-interval'9 days', v_stage_id, 'Hotfix deploy v1.4.1 to production', v_m1, 'Passed'),
  (NOW()-interval'9 days', NOW()-interval'9 days', v_stage_id, 'Verify notification badge on first login', v_m2, 'Passed');

  -- ── v1.5 — In Progress ─────────────────────────────────────────────────────
  v_pid := v_project_ids[1+floor(random()*v_nproj)::int];
  v_m1 := v_member_ids[1+floor(random()*v_nmem)::int];
  v_m2 := v_member_ids[1+floor(random()*v_nmem)::int];
  INSERT INTO releases (created_at, updated_at, project_id, name, description, status, target_date)
  VALUES (NOW()-interval'5 days', NOW(), v_pid,
    'v1.5.0', 'Audit log, WebSocket notifications, timezone support, performance improvements.', 'In Progress', '2026-06-20')
  RETURNING id INTO v_release_id;

  INSERT INTO release_stages (created_at, updated_at, release_id, name, status)
  VALUES (NOW()-interval'5 days', NOW(), v_release_id, 'Development', 'Active')
  RETURNING id INTO v_stage_id;
  INSERT INTO release_stories (created_at, updated_at, stage_id, title, assignee_id, status) VALUES
  (NOW()-interval'5 days', NOW()-interval'2 days', v_stage_id, 'Audit log model, migration and service', v_m2, 'In QA'),
  (NOW()-interval'4 days', NOW()-interval'1 day',  v_stage_id, 'Audit log admin UI', v_m2, 'Pending'),
  (NOW()-interval'4 days', NOW()-interval'2 days', v_stage_id, 'Timezone refactor (11 templates)', v_m1, 'Passed'),
  (NOW()-interval'3 days', NOW()-interval'1 day',  v_stage_id, 'Fix N+1 on sprint board and release detail', v_m1, 'Pending'),
  (NOW()-interval'2 days', NOW()-interval'2 days', v_stage_id, 'WebSocket notification spike', v_m1, 'Pending');

  INSERT INTO release_stages (created_at, updated_at, release_id, name, status)
  VALUES (NOW()-interval'1 day', NOW(), v_release_id, 'QA', 'Pending')
  RETURNING id INTO v_stage_id;
  INSERT INTO release_stories (created_at, updated_at, stage_id, title, assignee_id, status) VALUES
  (NOW()-interval'1 day', NOW()-interval'1 day', v_stage_id, 'Full regression suite', v_m2, 'Pending'),
  (NOW()-interval'1 day', NOW()-interval'1 day', v_stage_id, 'Performance regression check', v_m1, 'Pending');

  INSERT INTO release_stages (created_at, updated_at, release_id, name, status)
  VALUES (NOW(), NOW(), v_release_id, 'Production Deploy', 'Pending')
  RETURNING id INTO v_stage_id;
  INSERT INTO release_stories (created_at, updated_at, stage_id, title, assignee_id, status) VALUES
  (NOW(), NOW(), v_stage_id, 'Deploy v1.5.0 to production', v_m1, 'Pending'),
  (NOW(), NOW(), v_stage_id, 'Post-deploy monitoring — 24h watch', v_m2, 'Pending');

  RAISE NOTICE 'Done — 4 releases with stages and stories inserted.';
END $$;
