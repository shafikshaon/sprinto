package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sprinto/models"
	"sprinto/service"
)

type SprintsData struct {
	Meta   Meta
	Sprint models.Sprint
	Stats  models.SprintStats
}

type SprintHandler struct {
	svc service.SprintService
}

func NewSprintHandler(svc service.SprintService) *SprintHandler {
	return &SprintHandler{svc: svc}
}

func (h *SprintHandler) List(c *gin.Context) {
	sprint, err := h.svc.ActiveSprint()
	if err != nil {
		c.String(500, "DB error: %s", err.Error())
		return
	}
	sprintLabel := sprint.Name
	if sprint.StartDate != "" && sprint.EndDate != "" {
		sprintLabel += " · " + sprint.StartDate + " – " + sprint.EndDate
	}
	allProjects, activeProject := projectMeta(c)
	render(c, "sprints", SprintsData{
		Meta:   Meta{Title: "Sprint Board", CurrentPage: "sprints", SprintLabel: sprintLabel, AllProjects: allProjects, ActiveProject: activeProject},
		Sprint: sprint,
		Stats:  models.ComputeStats(sprint.Tasks),
	})
}

func (h *SprintHandler) CreateTask(c *gin.Context) {
	sprintID, _ := strconv.ParseUint(c.PostForm("sprint_id"), 10, 64)
	h.svc.AddTask(
		uint(sprintID),
		c.PostForm("title"),
		c.PostForm("assignee"),
		c.PostForm("status"),
		c.PostForm("priority"),
	)
	redirectTo(c, "/sprints")
}

func (h *SprintHandler) DeleteTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.RemoveTask(uint(id))
	redirectTo(c, "/sprints")
}

func (h *SprintHandler) UpdateProgress(c *gin.Context) {
	sprintID, _ := strconv.ParseUint(c.PostForm("sprint_id"), 10, 64)
	progress, _ := strconv.Atoi(c.PostForm("progress"))
	h.svc.UpdateProgress(uint(sprintID), progress)
	redirectTo(c, "/sprints")
}
