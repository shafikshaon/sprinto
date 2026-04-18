package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sprinto/models"
	"sprinto/service"
)

type SprintsData struct {
	Meta    Meta
	Sprint  models.Sprint
	Stats   models.SprintStats
	Members []models.TeamMember
}

type SprintTaskData struct {
	Meta    Meta
	Task    models.SprintTask
	Members []models.TeamMember
}

type SprintHandler struct {
	svc service.SprintService
}

func NewSprintHandler(svc service.SprintService) *SprintHandler {
	return &SprintHandler{svc: svc}
}

func (h *SprintHandler) List(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)
	sprint, err := h.svc.ActiveSprint(projectID)
	if err != nil {
		c.String(500, "DB error: %s", err.Error())
		return
	}
	sprintLabel := sprint.Name
	if sprint.StartDate != "" && sprint.EndDate != "" {
		sprintLabel += " · " + sprint.StartDate + " – " + sprint.EndDate
	}
	allProjects, activeProject := projectMeta(c)
	var members []models.TeamMember
	if activeProject != nil {
		members = activeProject.Members
	}
	render(c, "sprints", SprintsData{
		Meta:    Meta{Title: "Sprint Board", CurrentPage: "sprints", ActionLabel: "Add Task", SprintLabel: sprintLabel, AllProjects: allProjects, ActiveProject: activeProject},
		Sprint:  sprint,
		Stats:   models.ComputeStats(sprint.Tasks),
		Members: members,
	})
}

func (h *SprintHandler) CreateTask(c *gin.Context) {
	sprintID, _ := strconv.ParseUint(c.PostForm("sprint_id"), 10, 64)
	h.svc.AddTask(
		uint(sprintID),
		c.PostForm("title"),
		c.PostFormArray("assignees"),
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

func (h *SprintHandler) TaskDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	task, err := h.svc.TaskByID(uint(id))
	if err != nil {
		c.String(404, "Task not found")
		return
	}
	allProjects, activeProject := projectMeta(c)
	var members []models.TeamMember
	if activeProject != nil {
		members = activeProject.Members
	}
	render(c, "sprint_task", SprintTaskData{
		Meta:    Meta{Title: task.Title, CurrentPage: "sprints", AllProjects: allProjects, ActiveProject: activeProject},
		Task:    task,
		Members: members,
	})
}

func (h *SprintHandler) UpdateTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.UpdateTask(
		uint(id),
		c.PostForm("title"),
		c.PostFormArray("assignees"),
		c.PostForm("status"),
		c.PostForm("priority"),
	)
	redirectTo(c, "/sprints/tasks/"+c.Param("id"))
}

func (h *SprintHandler) AddComment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.AddComment(uint(id), c.PostForm("author"), c.PostForm("content"))
	redirectTo(c, "/sprints/tasks/"+c.Param("id"))
}

func (h *SprintHandler) DeleteComment(c *gin.Context) {
	commentID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	taskID := c.PostForm("task_id")
	h.svc.DeleteComment(uint(commentID))
	redirectTo(c, "/sprints/tasks/"+taskID)
}
