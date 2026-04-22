# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the application
go run main.go

# Build
go build ./...

# Tidy dependencies
go mod tidy
```

The app requires PostgreSQL. Set env vars or rely on defaults:
- `DATABASE_URL`: `host=localhost port=5432 dbname=sprinto user=postgres password=postgres sslmode=disable`
- `PORT`: `8080`
- `SESSION_SECRET`: HMAC key for session cookies (defaults to a dev fallback — change in prod)

Seed data is loaded by running `seed.sql` against the database:
```bash
psql $DATABASE_URL -f seed.sql
```
The seed requires at least one user and projects to exist first. Run the app once to register a user, then apply the seed.

## Architecture

3-layer: **Handler → Service → Repository**, wired in `main.go`.

```
main.go              Wire-up only — DB, repos, services, handlers, Gin routes
config/config.go     DATABASE_URL + PORT env vars
db/db.go             Connect, AutoMigrate (all models), legacy schema fixups
models/models.go     All GORM model structs in one file
repository/          One file per entity: interface + GORM implementation
service/             One file per entity: interface + implementation
handlers/            One file per entity + shared helpers in handlers.go
templates/           layout.html + one content template per page
store/               Deprecated — stub package only, kept for historical reasons
```

### Unified Task table

All task variants share the `tasks` table, differentiated by `category`:
- `"sprint"` — sprint board tasks; linked via `sprint_id`
- `"dev"` — dev backlog tasks; linked via `project_id`
- `"release"` — release checklist stories; linked via `release_stage_id`

Sprint task assignees use the `task_assignees` many2many junction (`task_id`, `team_member_id`). Release stories use a single `assignee_id` FK instead.

### Release data lives in Sprint

There is no separate releases table. Release metadata (`description`, `status`, `target_date`) are columns on the `sprints` table. Release stages are `release_stages` rows with `sprint_id` FK. The `release_id` column on `release_stages` is a legacy NOT-NULL column from the old schema — `db/db.go:Migrate()` drops the constraint at startup.

### Templates

- `layout.html` defines `{{define "layout"}}`. Every page template defines `{{define "content"}}`.
- Auth pages (`login.html`, `register.html`) skip the layout entirely and use `renderAuth()` which parses only the standalone template.
- Templates are **parsed per request** via `render(c, "page-name", data)` in `handlers/handlers.go` — `Gin.LoadHTMLGlob` is not used.
- Template function map (`funcMap` in `handlers/handlers.go`): `initials`, `statusClass`, `priorityClass`, `typeClass`, `urgencyClass`, `timeAgo`, `gregorianDate`, `hijriDate`, `bengaliDate`, `add`, `sub`, `wasEdited`.

### Context pipeline (every protected request)

1. `LoadUserMiddleware` — reads HMAC-signed session cookie (`sprinto_session`), sets `"current_user"` in Gin context.
2. `AuthRequiredMiddleware` — redirects to `/login` if no user.
3. `ProjectMiddleware` — loads all projects + resolves active project from `active_project` cookie; sets `"all_projects"` and `"active_project"` in context.
4. Handler calls `projectMeta(c)` to retrieve the above three values, then builds a `Meta` struct and calls `render()`.

### Meta struct

Passed to every template. Key fields:
- `CurrentPage` — matches nav link names for active-state highlighting in the sidebar
- `ActionLabel` — if non-empty, renders a violet "+ Label" button in the page header that opens `#add-dialog`
- `SprintLabel` — shown in the top bar (e.g. `"Sprint 14 · May 25 – Jun 6"`)
- `AllProjects` / `ActiveProject` / `CurrentUser` — populated from Gin context via `projectMeta(c)`

### Computed (non-persisted) fields

Populated in the service layer after DB load, tagged `gorm:"-"`:
- `Deadline.DueDate` (formatted string), `Deadline.DaysLeft` (int)
- `Meeting.Attendees` ([]string split from `AttendeeCSV`)

### Custom table names

| Model | Table |
|---|---|
| `StandupEntry` | `standups` |
| `ActionItem` | `meeting_action_items` |
| `Task` | `tasks` |
| `TaskComment` | `task_comments` |
| `ReleaseStage` | `release_stages` |
| `ReleaseSlackUpdate` | `release_slack_updates` |
| `TeamMember` | `team_members` |

### ID parsing

`gorm.Model` uses `uint`. Always parse path params with:
```go
id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
// then use uint(id)
```

### Forms and redirects

All mutating actions use `POST` forms with a 303 redirect back to the list page — no JSON API, no client-side fetch. Multi-select assignees use `c.PostFormArray("assignees")` → `parseUintArray()`.

### Dialog convention

Every list page that has a "+ Add" button uses `<dialog id="add-dialog">`. The `ActionLabel` field in `Meta` wires the header button to `document.getElementById('add-dialog').showModal()`. The global delete confirmation dialog (`confirm-delete-dialog`) is defined in `layout.html` and triggered by any button with `data-delete-action` / `data-delete-name` attributes.

### Design system

Tailwind CSS via CDN (no build step). Primary colour: `violet-600`. Key conventions:
- Edit buttons: `text-violet-600 bg-violet-50 border border-violet-200`
- Delete buttons: `text-red-600 bg-red-50 border border-red-200`
- Input focus: `focus:outline-none focus:border-violet-600 focus:ring-1 focus:ring-violet-600`
- Dialog header: `px-5 py-3.5 border-b border-gray-100`
- Dialog body: `px-5 py-4 space-y-3`

### Complete route list

| Method | Path | Handler |
|---|---|---|
| GET/POST | `/login`, `/register`, POST `/logout` | AuthHandler |
| GET | `/` | DashboardHandler |
| GET/POST | `/sprints` + sub-routes | SprintHandler |
| GET/POST | `/standups`, `/standups/pdf` | StandupHandler |
| GET/POST | `/deadlines` | DeadlineHandler |
| GET/POST | `/meetings` | MeetingHandler |
| GET/POST | `/devtasks`, `/devtasks/:id` | DevTaskHandler |
| GET/POST | `/projects`, `/switch-project` | ProjectHandler |
| GET/POST | `/team` | TeamHandler |
| GET/POST | `/notes`, `/notes/new`, `/notes/:id/edit` | StickyNoteHandler |
| GET/POST | `/slack` | SlackHandler |
