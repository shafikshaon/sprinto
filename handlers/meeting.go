package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sprinto/models"
	"sprinto/service"
)

type MeetingsData struct {
	Meta     Meta
	Meetings []models.Meeting
}

type MeetingHandler struct {
	svc service.MeetingService
}

func NewMeetingHandler(svc service.MeetingService) *MeetingHandler {
	return &MeetingHandler{svc: svc}
}

func (h *MeetingHandler) List(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)
	meetings, err := h.svc.All(projectID)
	if err != nil {
		c.String(500, "DB error: %s", err.Error())
		return
	}
	allProjects, activeProject := projectMeta(c)
	render(c, "meetings", MeetingsData{
		Meta:     Meta{Title: "Meeting Minutes", CurrentPage: "meetings", AllProjects: allProjects, ActiveProject: activeProject},
		Meetings: meetings,
	})
}

func (h *MeetingHandler) Create(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)
	h.svc.Add(
		c.PostForm("title"),
		c.PostForm("date"),
		c.PostForm("attendees"),
		c.PostForm("notes"),
		projectID,
	)
	redirectTo(c, "/meetings")
}

func (h *MeetingHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Remove(uint(id))
	redirectTo(c, "/meetings")
}
