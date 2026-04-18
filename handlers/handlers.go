// Package handlers contains Gin HTTP handlers.
// Each handler calls a service; no business logic or DB access here.
package handlers

import (
	"html/template"
	"log"
	"strings"

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
}

// projectMeta reads project context injected by ProjectMiddleware.
func projectMeta(c *gin.Context) ([]models.Project, *models.Project) {
	all, _ := c.Get("all_projects")
	active, _ := c.Get("active_project")
	projects, _ := all.([]models.Project)
	activePrj, _ := active.(*models.Project)
	return projects, activePrj
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
	"urgencyClass": func(days int) string {
		if days <= 3 {
			return "text-red-600 font-semibold"
		} else if days <= 7 {
			return "text-yellow-600 font-semibold"
		}
		return "text-gray-400"
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
