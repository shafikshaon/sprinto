# Sprinto

A self-hosted team productivity tool for engineering teams — covering sprint management, standups, deadlines, meetings, dev tasks, and release tracking in a single lightweight web app.

Built with Go, Gin, GORM, PostgreSQL, and Tailwind CSS (CDN). No frontend build step required.

---

## Features

- **Sprint Board** — manage active sprint tasks with status, priority, and assignees; track progress with a manual slider; view all past sprints
- **Release Tracker** — define release stages and checklist stories per sprint; track QA/pass/fail status and attach Slack update threads
- **Dev Tasks** — backlog of improvements, tech debt, and research items with comments and multi-assignee support
- **Standups** — daily standup log with filtering, pagination, and PDF export
- **Deadlines** — milestone tracker with urgency grouping (Critical / High / Medium / Low) and project association
- **Meetings** — meeting minutes with attendees, notes, and action items
- **Slack Threads** — saved Slack discussion summaries with tags and author attribution
- **Notes** — sticky notes with colour coding, pinning, Markdown preview, and live editor
- **Projects** — project management with member associations and an active-project switcher in the sidebar
- **Team** — team member directory linked to user accounts

---

## Requirements

- Go 1.21+
- PostgreSQL 14+

---

## Getting Started

### 1. Create the database

```bash
createdb sprinto
```

### 2. Configure environment

```bash
export DATABASE_URL="host=localhost port=5432 dbname=sprinto user=postgres password=postgres sslmode=disable"
export PORT=8080
export SESSION_SECRET="replace-with-a-random-secret-in-production"
```

Defaults are used if these are not set (`DATABASE_URL` above, port `8080`, a dev-only session secret).

### 3. Run

```bash
go run main.go
```

On first start, `AutoMigrate` creates all tables. Open `http://localhost:8080/register` to create your account.

### 4. Load seed data (optional)

After registering and creating at least one project:

```bash
psql $DATABASE_URL -f seed.sql
```

The seed file inserts projects, team members, standups, dev tasks, sprints, releases, deadlines, meetings, Slack threads, and notes — enough data to explore all features immediately.

---

## Build

```bash
go build ./...
```

---

## Project Structure

```
main.go           Wire-up: DB, repos, services, handlers, routes
config/           Environment variable loading
db/               Database connection and AutoMigrate
models/           All GORM model structs (single file)
repository/       Data access interfaces + GORM implementations
service/          Business logic interfaces + implementations
handlers/         Gin HTTP handlers + template rendering
templates/        Go html/template files (layout.html + one per page)
seed.sql          Sample data for local development
```

---

## Stack

| Layer | Technology |
|---|---|
| Web framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io) with PostgreSQL driver |
| Templates | Go `html/template` (server-side rendering) |
| Styling | [Tailwind CSS](https://tailwindcss.com) via CDN |
| Auth | HMAC-signed session cookie (`bcrypt` passwords) |
| PDF export | [go-pdf/fpdf](https://github.com/go-pdf/fpdf) |
