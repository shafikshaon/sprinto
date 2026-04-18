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

# Run tests (none yet)
go test ./...
```

The app requires a running PostgreSQL instance. Set `DATABASE_URL` and `PORT` env vars, or rely on defaults:
- `DATABASE_URL`: `host=localhost port=5432 dbname=sprinto user=postgres password=postgres sslmode=disable`
- `PORT`: `8080`

On first run with an empty database, seed data is inserted automatically.

## Architecture

3-layer architecture: **Handler → Service → Repository**

```
main.go              Wire-up only: DB, repos, services, handlers, Gin routes
config/config.go     Reads DATABASE_URL and PORT env vars
db/db.go             Connect, AutoMigrate (all models), IsEmpty, Seed
models/models.go     GORM model structs for all entities
repository/          Repository interfaces + GORM implementations (one file)
service/             Service interfaces + implementations (one file)
handlers/            Gin HTTP handlers (one file per entity + handlers.go)
templates/           html/template files: layout.html + one per page
```

### Key conventions

- **Templates**: `layout.html` defines `{{define "layout"}}`, each page defines `{{define "content"}}`. Templates are parsed per-request via `render()` in `handlers/handlers.go` — not using Gin's `LoadHTMLGlob`.
- **Meta struct** (`handlers/handlers.go`): passed to every template; controls sidebar active state, page title, optional top-bar action button, and `SprintLabel` for the top bar.
- **Computed fields**: `Deadline.DueDate`, `Deadline.DaysLeft`, and `Meeting.Attendees` are tagged `gorm:"-"` and populated in the service layer after DB load.
- **IDs**: GORM's `gorm.Model` uses `uint` — parse path params with `strconv.ParseUint(..., 10, 64)` then cast to `uint`.
- **Custom table names**: `StandupEntry` → `standups`, `ActionItem` → `meeting_action_items`, `DevTask` → `dev_tasks` (via `TableName()` methods).
- **Template helpers** (defined in `handlers/handlers.go` `funcMap`): `initials`, `statusClass`, `priorityClass`, `typeClass`, `urgencyClass`.

### Pages and routes

| Route | Handler | Template |
|---|---|---|
| `GET /` | DashboardHandler | dashboard.html |
| `GET /sprints` | SprintHandler | sprints.html |
| `GET /standups?date=` | StandupHandler | standups.html |
| `GET /deadlines` | DeadlineHandler | deadlines.html |
| `GET /meetings` | MeetingHandler | meetings.html |
| `GET /devtasks` | DevTaskHandler | devtasks.html |

All create/delete actions use `POST` forms (no JS) with 303 redirect back to the list page.
