package handlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"sprinto/models"
	"sprinto/service"
)

type StandupsData struct {
	Meta     Meta
	Standups []models.StandupEntry
	Today    string
}

type StandupHandler struct {
	svc service.StandupService
}

func NewStandupHandler(svc service.StandupService) *StandupHandler {
	return &StandupHandler{svc: svc}
}

func (h *StandupHandler) List(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)
	entries, _ := h.svc.All(projectID)
	allProjects, activeProject, currentUser := projectMeta(c)
	render(c, "standups", StandupsData{
		Meta:     Meta{Title: "Daily Standups", CurrentPage: "standups", ActionLabel: "Add Standup", AllProjects: allProjects, ActiveProject: activeProject, CurrentUser: currentUser},
		Standups: entries,
		Today:    time.Now().Format("2006-01-02"),
	})
}

func (h *StandupHandler) Create(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)
	h.svc.Add(
		c.PostForm("date"),
		c.PostForm("summary"),
		c.PostForm("dependencies"),
		c.PostForm("issues"),
		c.PostForm("action_items"),
		projectID,
	)
	redirectTo(c, "/standups")
}

func (h *StandupHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Update(
		uint(id),
		c.PostForm("summary"),
		c.PostForm("dependencies"),
		c.PostForm("issues"),
		c.PostForm("action_items"),
	)
	redirectTo(c, "/standups")
}

func (h *StandupHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Remove(uint(id))
	redirectTo(c, "/standups")
}
