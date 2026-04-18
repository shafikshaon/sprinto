package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sprinto/service"
)

type SlackThreadsData struct {
	Meta    Meta
	Threads []threadRow
	AllTags []string
	Tag     string
}

type threadRow struct {
	ID          uint
	MessageLink string
	Topic       string
	Summary     string
	Tags        []string
	TagCSV      string
	Author      string
}

type SlackHandler struct {
	svc service.SlackThreadService
}

func NewSlackHandler(svc service.SlackThreadService) *SlackHandler {
	return &SlackHandler{svc: svc}
}

func (h *SlackHandler) List(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)
	tag := c.Query("tag")

	threads, err := h.svc.All(projectID, tag)
	if err != nil {
		c.String(500, "DB error: %s", err.Error())
		return
	}
	allTags, _ := h.svc.AllTags(projectID)

	rows := make([]threadRow, len(threads))
	for i, t := range threads {
		rows[i] = threadRow{
			ID:          t.ID,
			MessageLink: t.MessageLink,
			Topic:       t.Topic,
			Summary:     t.Summary,
			Tags:        t.Tags,
			TagCSV:      t.TagCSV,
			Author:      t.Author,
		}
	}

	allProjects, activeProject, currentUser := projectMeta(c)
	render(c, "slack", SlackThreadsData{
		Meta:    Meta{Title: "Slack Threads", CurrentPage: "slack", ActionLabel: "Add Thread", AllProjects: allProjects, ActiveProject: activeProject, CurrentUser: currentUser},
		Threads: rows,
		AllTags: allTags,
		Tag:     tag,
	})
}

func (h *SlackHandler) Create(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)
	h.svc.Create(
		c.PostForm("message_link"),
		c.PostForm("topic"),
		c.PostForm("summary"),
		c.PostForm("tags"),
		c.PostForm("author"),
		projectID,
	)
	redirectTo(c, "/slack")
}

func (h *SlackHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Update(
		uint(id),
		c.PostForm("message_link"),
		c.PostForm("topic"),
		c.PostForm("summary"),
		c.PostForm("tags"),
		c.PostForm("author"),
	)
	redirectTo(c, "/slack")
}

func (h *SlackHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Delete(uint(id))
	redirectTo(c, "/slack")
}
