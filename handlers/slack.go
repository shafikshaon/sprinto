package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sprinto/models"
	"sprinto/service"
)

type SlackThreadsData struct {
	Meta    Meta
	Threads []threadRow
	AllTags []string
	Tag     string
	Members []models.TeamMember
}

type threadRow struct {
	ID          uint
	MessageLink string
	Topic       string
	Summary     string
	Tags        []string
	TagCSV      string
	AuthorID    *uint
	AuthorName  string
}

type SlackHandler struct {
	svc service.SlackThreadService
}

func NewSlackHandler(svc service.SlackThreadService) *SlackHandler {
	return &SlackHandler{svc: svc}
}

func (h *SlackHandler) List(c *gin.Context) {
	tag := c.Query("tag")

	threads, err := h.svc.All(tag)
	if err != nil {
		c.String(500, "DB error: %s", err.Error())
		return
	}
	allTags, _ := h.svc.AllTags()

	rows := make([]threadRow, len(threads))
	for i, t := range threads {
		authorName := ""
		if t.Author != nil {
			authorName = t.Author.Name
		}
		rows[i] = threadRow{
			ID:          t.ID,
			MessageLink: t.MessageLink,
			Topic:       t.Topic,
			Summary:     t.Summary,
			Tags:        t.Tags,
			TagCSV:      t.TagCSV,
			AuthorID:    t.AuthorID,
			AuthorName:  authorName,
		}
	}

	allProjects, activeProject, currentUser := projectMeta(c)
	var members []models.TeamMember
	if activeProject != nil {
		members = activeProject.Members
	}
	render(c, "slack", SlackThreadsData{
		Meta:    Meta{Title: "Slack Threads", CurrentPage: "slack", ActionLabel: "Add Thread", AllProjects: allProjects, ActiveProject: activeProject, CurrentUser: currentUser},
		Threads: rows,
		AllTags: allTags,
		Tag:     tag,
		Members: members,
	})
}

func (h *SlackHandler) Create(c *gin.Context) {
	var authorID *uint
	if aid, err := strconv.ParseUint(c.PostForm("author"), 10, 64); err == nil && aid > 0 {
		v := uint(aid)
		authorID = &v
	}
	h.svc.Create(
		c.PostForm("message_link"),
		c.PostForm("topic"),
		c.PostForm("summary"),
		c.PostForm("tags"),
		authorID,
	)
	redirectTo(c, "/slack")
}

func (h *SlackHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var authorID *uint
	if aid, err := strconv.ParseUint(c.PostForm("author"), 10, 64); err == nil && aid > 0 {
		v := uint(aid)
		authorID = &v
	}
	h.svc.Update(
		uint(id),
		c.PostForm("message_link"),
		c.PostForm("topic"),
		c.PostForm("summary"),
		c.PostForm("tags"),
		authorID,
	)
	redirectTo(c, "/slack")
}

func (h *SlackHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Delete(uint(id))
	redirectTo(c, "/slack")
}
