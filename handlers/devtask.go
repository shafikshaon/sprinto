package handlers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"

	"sprinto/models"
	"sprinto/repository"
	"sprinto/service"
)

const devTaskPerPage = 15

type DevTasksData struct {
	Meta     Meta
	DevTasks []models.DevTask
	Counts   map[string]int
	Members  []models.TeamMember
}

type DevTaskDetailData struct {
	Meta    Meta
	Task    models.DevTask
	Members []models.TeamMember
}

type DevTaskHandler struct {
	svc service.DevTaskService
}

func NewDevTaskHandler(svc service.DevTaskService) *DevTaskHandler {
	return &DevTaskHandler{svc: svc}
}

func (h *DevTaskHandler) List(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	f := repository.DevTaskFilter{
		Search:   c.Query("search"),
		Type:     c.Query("type"),
		Status:   c.Query("status"),
		Priority: c.Query("priority"),
	}

	tasks, total, err := h.svc.All(projectID, f, page, devTaskPerPage)
	if err != nil {
		c.String(500, "DB error: %s", err.Error())
		return
	}
	counts, _ := h.svc.OpenCountsByType(projectID)

	totalPages := int(math.Ceil(float64(total) / float64(devTaskPerPage)))
	if totalPages < 1 {
		totalPages = 1
	}
	pages := make([]int, totalPages)
	for i := range pages {
		pages[i] = i + 1
	}

	allProjects, activeProject, currentUser := projectMeta(c)
	var members []models.TeamMember
	if activeProject != nil {
		members = activeProject.Members
	}
	render(c, "devtasks", map[string]interface{}{
		"Meta":       Meta{Title: "Dev Tasks & Improvements", CurrentPage: "devtasks", ActionLabel: "Add Task", AllProjects: allProjects, ActiveProject: activeProject, CurrentUser: currentUser},
		"DevTasks":   tasks,
		"Counts":     counts,
		"Members":    members,
		"Filter":     f,
		"Page":       page,
		"TotalPages": totalPages,
		"Total":      total,
		"Pages":      pages,
	})
}

func (h *DevTaskHandler) Create(c *gin.Context) {
	pid, _ := strconv.ParseUint(c.PostForm("project_id"), 10, 64)
	h.svc.Add(
		c.PostForm("title"),
		c.PostForm("type"),
		c.PostFormArray("assignees"),
		c.PostForm("status"),
		c.PostForm("priority"),
		uint(pid),
	)
	redirectTo(c, "/devtasks")
}

func (h *DevTaskHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Remove(uint(id))
	redirectTo(c, "/devtasks")
}

func (h *DevTaskHandler) Detail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	task, err := h.svc.ByID(uint(id))
	if err != nil {
		c.String(404, "Task not found")
		return
	}
	allProjects, activeProject, currentUser := projectMeta(c)
	var members []models.TeamMember
	if activeProject != nil {
		members = activeProject.Members
	}
	render(c, "devtask_detail", DevTaskDetailData{
		Meta:    Meta{Title: task.Title, CurrentPage: "devtasks", AllProjects: allProjects, ActiveProject: activeProject, CurrentUser: currentUser},
		Task:    task,
		Members: members,
	})
}

func (h *DevTaskHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Update(
		uint(id),
		c.PostForm("title"),
		c.PostForm("type"),
		c.PostFormArray("assignees"),
		c.PostForm("status"),
		c.PostForm("priority"),
	)
	redirectTo(c, "/devtasks/"+c.Param("id"))
}

func (h *DevTaskHandler) AddComment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.AddComment(uint(id), c.PostForm("author"), c.PostForm("content"))
	redirectTo(c, "/devtasks/"+c.Param("id"))
}

func (h *DevTaskHandler) DeleteComment(c *gin.Context) {
	commentID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	taskID := c.PostForm("task_id")
	h.svc.DeleteComment(uint(commentID))
	redirectTo(c, "/devtasks/"+taskID)
}
