package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sprinto/models"
	"sprinto/service"
)

type DeadlineUrgency struct {
	Critical int // DaysLeft <= 3
	High     int // DaysLeft 4–14
	Medium   int // DaysLeft 15–30
	Low      int // DaysLeft > 30
}

type DeadlinesData struct {
	Meta      Meta
	Deadlines []models.Deadline
	Urgency   DeadlineUrgency
}

type DeadlineHandler struct {
	svc service.DeadlineService
}

func NewDeadlineHandler(svc service.DeadlineService) *DeadlineHandler {
	return &DeadlineHandler{svc: svc}
}

func (h *DeadlineHandler) List(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)
	deadlines, err := h.svc.All(projectID)
	if err != nil {
		c.String(500, "DB error: %s", err.Error())
		return
	}
	var urgency DeadlineUrgency
	for _, d := range deadlines {
		switch {
		case d.DaysLeft <= 3:
			urgency.Critical++
		case d.DaysLeft <= 14:
			urgency.High++
		case d.DaysLeft <= 30:
			urgency.Medium++
		default:
			urgency.Low++
		}
	}

	allProjects, activeProject := projectMeta(c)
	render(c, "deadlines", DeadlinesData{
		Meta:      Meta{Title: "Deadlines", CurrentPage: "deadlines", AllProjects: allProjects, ActiveProject: activeProject},
		Deadlines: deadlines,
		Urgency:   urgency,
	})
}

func (h *DeadlineHandler) Create(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)
	h.svc.Add(
		c.PostForm("title"),
		c.PostForm("project"),
		c.PostForm("due_date"),
		c.PostForm("priority"),
		projectID,
	)
	redirectTo(c, "/deadlines")
}

func (h *DeadlineHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Remove(uint(id))
	redirectTo(c, "/deadlines")
}
