// Package handlers contains Gin HTTP handlers.
// Each handler calls a service; no business logic or DB access here.
package handlers

import (
	"fmt"
	"html/template"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"sprinto/models"
)

// Meta holds data common to every page (sidebar state, top bar).
type Meta struct {
	Title         string
	CurrentPage   string
	ActionLabel   string
	ActionHref    string
	SprintLabel   string // e.g. "Sprint 12 · Apr 14 – 28"
	AllProjects   []models.Project
	ActiveProject *models.Project
	CurrentUser   *models.User
}

// projectMeta reads project context injected by ProjectMiddleware.
func projectMeta(c *gin.Context) ([]models.Project, *models.Project, *models.User) {
	all, _ := c.Get("all_projects")
	active, _ := c.Get("active_project")
	cu, _ := c.Get("current_user")
	projects, _ := all.([]models.Project)
	activePrj, _ := active.(*models.Project)
	currentUser, _ := cu.(*models.User)
	return projects, activePrj, currentUser
}

// activeProjectIDFromCtx returns the active project's ID (0 if none selected).
func activeProjectIDFromCtx(c *gin.Context) uint {
	_, active, _ := projectMeta(c)
	if active != nil {
		return active.ID
	}
	return 0
}

var funcMap = template.FuncMap{
	"initials": func(name string) string {
		parts := strings.Fields(name)
		r := ""
		for _, p := range parts {
			if len(p) > 0 {
				r += string([]rune(p)[0])
			}
		}
		if len([]rune(r)) > 2 {
			return string([]rune(r)[:2])
		}
		return r
	},
	"statusClass": func(status string) string {
		switch status {
		case "Done", "On Track", "Released", "Passed":
			return "bg-green-50 text-green-700"
		case "In Progress", "Active", "In QA":
			return "bg-blue-50 text-blue-700"
		case "Blocked", "Rolled Back", "Failed":
			return "bg-red-50 text-red-700"
		case "At Risk":
			return "bg-yellow-50 text-yellow-700"
		default:
			return "bg-gray-100 text-gray-600"
		}
	},
	"priorityClass": func(p string) string {
		switch p {
		case "Critical":
			return "bg-red-50 text-red-700"
		case "High":
			return "bg-orange-50 text-orange-700"
		case "Medium":
			return "bg-yellow-50 text-yellow-700"
		default:
			return "bg-gray-100 text-gray-500"
		}
	},
	"typeClass": func(t string) string {
		switch t {
		case "Tech Debt":
			return "bg-orange-50 text-orange-700"
		case "Improvement":
			return "bg-blue-50 text-blue-700"
		case "Research":
			return "bg-purple-50 text-purple-700"
		default:
			return "bg-gray-100 text-gray-600"
		}
	},
	"gregorianDate": func() string { return gregorianDate(time.Now()) },
	"hijriDate":     func() string { return hijriDate(time.Now()) },
	"bengaliDate":   func() string { return bengaliDate(time.Now()) },
	"timeAgo": func(t time.Time) string {
		d := time.Since(t)
		switch {
		case d < time.Minute:
			return "just now"
		case d < time.Hour:
			m := int(d.Minutes())
			if m == 1 {
				return "1 minute ago"
			}
			return fmt.Sprintf("%d minutes ago", m)
		case d < 24*time.Hour:
			h := int(d.Hours())
			if h == 1 {
				return "1 hour ago"
			}
			return fmt.Sprintf("%d hours ago", h)
		case d < 7*24*time.Hour:
			days := int(d.Hours() / 24)
			if days == 1 {
				return "1 day ago"
			}
			return fmt.Sprintf("%d days ago", days)
		default:
			return t.Format("Jan 2, 2006")
		}
	},
	"urgencyClass": func(days int) string {
		if days <= 3 {
			return "text-red-600 font-semibold"
		} else if days <= 7 {
			return "text-yellow-600 font-semibold"
		}
		return "text-gray-400"
	},
	"add":     func(a, b int) int { return a + b },
	"sub":     func(a, b int) int { return a - b },
	"wasEdited": func(created, updated time.Time) bool {
		return updated.Unix()-created.Unix() > 1
	},
}

// render parses layout + page template and writes the response.
func render(c *gin.Context, page string, data interface{}) {
	t, err := template.New("").Funcs(funcMap).ParseFiles(
		"templates/layout.html",
		"templates/"+page+".html",
	)
	if err != nil {
		c.String(500, "Template error: %s", err.Error())
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(200)
	if err := t.ExecuteTemplate(c.Writer, "layout", data); err != nil {
		log.Printf("render %s: %v", page, err)
	}
}

func redirectTo(c *gin.Context, path string) {
	c.Redirect(303, path)
}

// parseUintArray converts a slice of string IDs to []uint, skipping invalid/zero values.
func parseUintArray(strs []string) []uint {
	var result []uint
	for _, s := range strs {
		if v, err := strconv.ParseUint(s, 10, 64); err == nil && v > 0 {
			result = append(result, uint(v))
		}
	}
	return result
}
