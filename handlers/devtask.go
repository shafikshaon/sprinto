package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sprinto/models"
	"sprinto/service"
)

type DevTasksData struct {
	Meta     Meta
	DevTasks []models.DevTask
	Counts   map[string]int
}

type DevTaskHandler struct {
	svc service.DevTaskService
}

func NewDevTaskHandler(svc service.DevTaskService) *DevTaskHandler {
	return &DevTaskHandler{svc: svc}
}

func (h *DevTaskHandler) List(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)
	tasks, err := h.svc.All(projectID)
	if err != nil {
		c.String(500, "DB error: %s", err.Error())
		return
	}
	counts, _ := h.svc.OpenCountsByType(projectID)
	allProjects, activeProject := projectMeta(c)
	render(c, "devtasks", DevTasksData{
		Meta:     Meta{Title: "Dev Tasks & Improvements", CurrentPage: "devtasks", ActionLabel: "Add Task", AllProjects: allProjects, ActiveProject: activeProject},
		DevTasks: tasks,
		Counts:   counts,
	})
}

func (h *DevTaskHandler) Create(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)
	h.svc.Add(
		c.PostForm("title"),
		c.PostForm("type"),
		c.PostForm("assignee"),
		c.PostForm("status"),
		c.PostForm("priority"),
		projectID,
	)
	redirectTo(c, "/devtasks")
}

func (h *DevTaskHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Remove(uint(id))
	redirectTo(c, "/devtasks")
}
