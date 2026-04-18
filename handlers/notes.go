package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sprinto/service"
)

type NotesData struct {
	Meta  Meta
	Notes interface{}
}

type StickyNoteHandler struct {
	svc service.StickyNoteService
}

func NewStickyNoteHandler(svc service.StickyNoteService) *StickyNoteHandler {
	return &StickyNoteHandler{svc: svc}
}

func (h *StickyNoteHandler) List(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)
	filter := c.Query("filter")
	if filter == "" {
		filter = "all"
	}
	notes, err := h.svc.All(projectID, filter)
	if err != nil {
		c.String(500, "DB error: %s", err.Error())
		return
	}
	allProjects, activeProject, currentUser := projectMeta(c)
	render(c, "notes", gin.H{
		"Meta":   Meta{Title: "Sticky Notes", CurrentPage: "notes", ActionLabel: "New Note", AllProjects: allProjects, ActiveProject: activeProject, CurrentUser: currentUser},
		"Notes":  notes,
		"Filter": filter,
	})
}

func (h *StickyNoteHandler) New(c *gin.Context) {
	allProjects, activeProject, currentUser := projectMeta(c)
	render(c, "notes_new", gin.H{
		"Meta": Meta{Title: "New Note", CurrentPage: "notes", AllProjects: allProjects, ActiveProject: activeProject, CurrentUser: currentUser},
	})
}

func (h *StickyNoteHandler) EditPage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	note, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.Redirect(302, "/notes")
		return
	}
	allProjects, activeProject, currentUser := projectMeta(c)
	render(c, "notes_edit", gin.H{
		"Meta": Meta{Title: "Edit Note", CurrentPage: "notes", AllProjects: allProjects, ActiveProject: activeProject, CurrentUser: currentUser},
		"Note": note,
	})
}

func (h *StickyNoteHandler) Create(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)
	h.svc.Create(
		c.PostForm("title"),
		c.PostForm("content"),
		c.PostForm("color"),
		projectID,
	)
	redirectTo(c, "/notes")
}

func (h *StickyNoteHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Update(
		uint(id),
		c.PostForm("title"),
		c.PostForm("content"),
		c.PostForm("color"),
	)
	redirectTo(c, "/notes")
}

func (h *StickyNoteHandler) TogglePin(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	pinned := c.PostForm("pinned") == "true"
	h.svc.TogglePin(uint(id), pinned)
	redirect := c.PostForm("redirect")
	if redirect == "" {
		redirect = "/notes"
	}
	redirectTo(c, redirect)
}

func (h *StickyNoteHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Delete(uint(id))
	redirectTo(c, "/notes")
}
