package handlers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"

	"sprinto/models"
	"sprinto/repository"
	"sprinto/service"
)

const meetingPerPage = 10

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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	projectID, _ := strconv.ParseUint(c.Query("project"), 10, 64)
	f := repository.MeetingFilter{
		Search:   c.Query("search"),
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Project:  uint(projectID),
	}

	meetings, total, err := h.svc.All(f, page, meetingPerPage)
	if err != nil {
		c.String(500, "DB error: %s", err.Error())
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(meetingPerPage)))
	if totalPages < 1 {
		totalPages = 1
	}
	pages := make([]int, totalPages)
	for i := range pages {
		pages[i] = i + 1
	}

	allProjects, activeProject, currentUser := projectMeta(c)
	render(c, "meetings", map[string]interface{}{
		"Meta":       Meta{Title: "Meeting Minutes", CurrentPage: "meetings", ActionLabel: "New Meeting", AllProjects: allProjects, ActiveProject: activeProject, CurrentUser: currentUser},
		"Meetings":   meetings,
		"Filter":     f,
		"Page":       page,
		"TotalPages": totalPages,
		"Total":      total,
		"Pages":      pages,
	})
}

func (h *MeetingHandler) Create(c *gin.Context) {
	pid, _ := strconv.ParseUint(c.PostForm("project_id"), 10, 64)
	h.svc.Add(
		c.PostForm("title"),
		c.PostForm("date"),
		c.PostForm("attendees"),
		c.PostForm("notes"),
		uint(pid),
	)
	redirectTo(c, "/meetings")
}

func (h *MeetingHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	pid, _ := strconv.ParseUint(c.PostForm("project_id"), 10, 64)
	h.svc.Update(
		uint(id),
		uint(pid),
		c.PostForm("title"),
		c.PostForm("date"),
		c.PostForm("attendees"),
		c.PostForm("notes"),
	)
	redirectTo(c, "/meetings")
}

func (h *MeetingHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Remove(uint(id))
	redirectTo(c, "/meetings")
}
