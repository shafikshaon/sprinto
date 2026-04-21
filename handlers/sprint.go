package handlers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"sprinto/models"
	"sprinto/service"
)

type SprintsData struct {
	Meta       Meta
	Sprint     models.Sprint
	AllSprints []models.Sprint
	Stats      models.SprintStats
	Members    []models.TeamMember
}

type SprintTaskData struct {
	Meta    Meta
	Task    models.Task
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
	allSprints, err := h.svc.ListSprints(projectID)
	if err != nil {
		c.String(500, "DB error: %s", err.Error())
		return
	}
	sprintLabel := sprint.Name
	if sprint.StartDate != "" && sprint.EndDate != "" {
		sprintLabel += " · " + sprint.StartDate + " – " + sprint.EndDate
	}
	allProjects, activeProject, currentUser := projectMeta(c)
	var members []models.TeamMember
	if activeProject != nil {
		members = activeProject.Members
	}
	render(c, "sprints", SprintsData{
		Meta:       Meta{Title: "Sprint Board", CurrentPage: "sprints", ActionLabel: "New Sprint", SprintLabel: sprintLabel, AllProjects: allProjects, ActiveProject: activeProject, CurrentUser: currentUser},
		Sprint:     sprint,
		AllSprints: allSprints,
		Stats:      models.ComputeStats(sprint.Tasks),
		Members:    members,
	})
}

// ── Sprint Tasks ──────────────────────────────────────────────────────────────

func (h *SprintHandler) CreateTask(c *gin.Context) {
	sprintID, _ := strconv.ParseUint(c.PostForm("sprint_id"), 10, 64)
	h.svc.AddTask(
		uint(sprintID),
		c.PostForm("title"),
		parseUintArray(c.PostFormArray("assignees")),
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
	allProjects, activeProject, currentUser := projectMeta(c)
	var members []models.TeamMember
	if activeProject != nil {
		members = activeProject.Members
	}
	render(c, "sprint_task", SprintTaskData{
		Meta:    Meta{Title: task.Title, CurrentPage: "sprints", AllProjects: allProjects, ActiveProject: activeProject, CurrentUser: currentUser},
		Task:    task,
		Members: members,
	})
}

func (h *SprintHandler) UpdateTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.UpdateTask(
		uint(id),
		c.PostForm("title"),
		parseUintArray(c.PostFormArray("assignees")),
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

// ── Sprint Management ─────────────────────────────────────────────────────────

func (h *SprintHandler) CreateSprint(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)
	h.svc.CreateSprint(
		projectID,
		c.PostForm("name"),
		c.PostForm("goal"),
		c.PostForm("start_date"),
		c.PostForm("end_date"),
	)
	redirectTo(c, "/sprints")
}

func (h *SprintHandler) UpdateSprint(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.UpdateSprint(
		uint(id),
		c.PostForm("name"),
		c.PostForm("goal"),
		c.PostForm("start_date"),
		c.PostForm("end_date"),
	)
	redirectTo(c, "/sprints")
}

func (h *SprintHandler) DeleteSprint(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.DeleteSprint(uint(id))
	redirectTo(c, "/sprints")
}

func (h *SprintHandler) ActivateSprint(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	projectID := activeProjectIDFromCtx(c)
	h.svc.ActivateSprint(uint(id), projectID)
	redirectTo(c, "/sprints")
}

// ── Release (merged) ──────────────────────────────────────────────────────────

func (h *SprintHandler) UpdateRelease(c *gin.Context) {
	sprintID, _ := strconv.ParseUint(c.PostForm("sprint_id"), 10, 64)
	h.svc.UpdateRelease(uint(sprintID), c.PostForm("description"), c.PostForm("status"), c.PostForm("target_date"))
	redirectTo(c, "/sprints")
}

func (h *SprintHandler) CreateStage(c *gin.Context) {
	sprintID, _ := strconv.ParseUint(c.PostForm("sprint_id"), 10, 64)
	h.svc.AddStage(uint(sprintID), c.PostForm("name"), c.PostForm("status"))
	redirectTo(c, "/sprints")
}

func (h *SprintHandler) DeleteStage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.DeleteStage(uint(id))
	redirectTo(c, "/sprints")
}

func (h *SprintHandler) UpdateStageStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.UpdateStageStatus(uint(id), c.PostForm("status"))
	redirectTo(c, "/sprints")
}

func (h *SprintHandler) CreateStory(c *gin.Context) {
	stageID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var assigneeID *uint
	if aid, err := strconv.ParseUint(c.PostForm("assignee"), 10, 64); err == nil && aid > 0 {
		v := uint(aid)
		assigneeID = &v
	}
	h.svc.AddStory(uint(stageID), c.PostForm("title"), assigneeID)
	redirectTo(c, "/sprints")
}

func (h *SprintHandler) DeleteStory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.DeleteStory(uint(id))
	redirectTo(c, "/sprints")
}

func (h *SprintHandler) UpdateStoryStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.UpdateStoryStatus(uint(id), c.PostForm("status"))
	redirectTo(c, "/sprints")
}

func (h *SprintHandler) UpdateStory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var assigneeID *uint
	if aid, err := strconv.ParseUint(c.PostForm("assignee"), 10, 64); err == nil && aid > 0 {
		v := uint(aid)
		assigneeID = &v
	}
	h.svc.UpdateStory(uint(id), c.PostForm("title"), assigneeID)
	redirectTo(c, "/sprints")
}

func (h *SprintHandler) CreateSlackUpdate(c *gin.Context) {
	stageID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.AddSlackUpdate(
		uint(stageID),
		c.PostForm("channel"),
		c.PostForm("message"),
		fmt.Sprintf("%s", c.PostForm("author")),
	)
	redirectTo(c, "/sprints")
}

func (h *SprintHandler) DeleteSlackUpdate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.DeleteSlackUpdate(uint(id))
	redirectTo(c, "/sprints")
}
