package handlers

import (
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"sprinto/repository"
	"sprinto/service"
)

const standupPerPage = 14

type StandupHandler struct {
	svc service.StandupService
}

func NewStandupHandler(svc service.StandupService) *StandupHandler {
	return &StandupHandler{svc: svc}
}

func (h *StandupHandler) List(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	f := repository.StandupFilter{
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Search:   c.Query("search"),
	}

	entries, total, _ := h.svc.All(projectID, f, page, standupPerPage)
	totalPages := int(math.Ceil(float64(total) / float64(standupPerPage)))
	if totalPages < 1 {
		totalPages = 1
	}

	pages := make([]int, totalPages)
	for i := range pages {
		pages[i] = i + 1
	}

	allProjects, activeProject, currentUser := projectMeta(c)
	render(c, "standups", map[string]interface{}{
		"Meta":       Meta{Title: "Daily Standups", CurrentPage: "standups", ActionLabel: "Add Standup", AllProjects: allProjects, ActiveProject: activeProject, CurrentUser: currentUser},
		"Entries":    entries,
		"Today":      time.Now().Format("2006-01-02"),
		"Filter":     f,
		"Page":       page,
		"TotalPages": totalPages,
		"Total":      total,
		"Pages":      pages,
	})
}

func (h *StandupHandler) Create(c *gin.Context) {
	pid, _ := strconv.ParseUint(c.PostForm("project_id"), 10, 64)
	h.svc.Add(
		c.PostForm("date"),
		c.PostForm("summary"),
		c.PostForm("dependencies"),
		c.PostForm("issues"),
		c.PostForm("action_items"),
		uint(pid),
	)
	redirectTo(c, "/standups")
}

func (h *StandupHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	pid, _ := strconv.ParseUint(c.PostForm("project_id"), 10, 64)
	h.svc.Update(
		uint(id),
		c.PostForm("summary"),
		c.PostForm("dependencies"),
		c.PostForm("issues"),
		c.PostForm("action_items"),
		uint(pid),
	)
	redirectTo(c, "/standups")
}

func (h *StandupHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Remove(uint(id))
	redirectTo(c, "/standups")
}
