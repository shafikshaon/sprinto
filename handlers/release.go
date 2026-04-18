package handlers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"sprinto/models"
	"sprinto/service"
)

type ReleasesData struct {
	Meta     Meta
	Releases []models.Release
}

type ReleaseDetailData struct {
	Meta    Meta
	Release models.Release
}

type ReleaseHandler struct {
	svc service.ReleaseService
}

func NewReleaseHandler(svc service.ReleaseService) *ReleaseHandler {
	return &ReleaseHandler{svc: svc}
}

func (h *ReleaseHandler) List(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)
	releases, err := h.svc.All(projectID)
	if err != nil {
		c.String(500, "DB error: %s", err.Error())
		return
	}
	allProjects, activeProject, currentUser := projectMeta(c)
	render(c, "releases", ReleasesData{
		Meta:     Meta{Title: "Releases", CurrentPage: "releases", ActionLabel: "New Release", AllProjects: allProjects, ActiveProject: activeProject, CurrentUser: currentUser},
		Releases: releases,
	})
}

func (h *ReleaseHandler) Create(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)
	h.svc.Create(
		c.PostForm("name"),
		c.PostForm("description"),
		c.PostForm("status"),
		c.PostForm("target_date"),
		projectID,
	)
	redirectTo(c, "/releases")
}

func (h *ReleaseHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Delete(uint(id))
	redirectTo(c, "/releases")
}

func (h *ReleaseHandler) Detail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	release, err := h.svc.ByID(uint(id))
	if err != nil {
		c.String(404, "Release not found")
		return
	}
	allProjects, activeProject, currentUser := projectMeta(c)
	render(c, "release_detail", ReleaseDetailData{
		Meta:    Meta{Title: release.Name, CurrentPage: "releases", ActionLabel: "Add Stage", AllProjects: allProjects, ActiveProject: activeProject, CurrentUser: currentUser},
		Release: release,
	})
}

func (h *ReleaseHandler) CreateStage(c *gin.Context) {
	releaseID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.AddStage(uint(releaseID), c.PostForm("name"), c.PostForm("status"))
	redirectTo(c, fmt.Sprintf("/releases/%d", releaseID))
}

func (h *ReleaseHandler) DeleteStage(c *gin.Context) {
	stageID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	releaseID, _ := strconv.ParseUint(c.PostForm("release_id"), 10, 64)
	h.svc.DeleteStage(uint(stageID))
	redirectTo(c, fmt.Sprintf("/releases/%d", releaseID))
}

func (h *ReleaseHandler) UpdateStageStatus(c *gin.Context) {
	stageID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	releaseID, _ := strconv.ParseUint(c.PostForm("release_id"), 10, 64)
	h.svc.UpdateStageStatus(uint(stageID), c.PostForm("status"))
	redirectTo(c, fmt.Sprintf("/releases/%d", releaseID))
}

func (h *ReleaseHandler) CreateStory(c *gin.Context) {
	stageID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	releaseID, _ := strconv.ParseUint(c.PostForm("release_id"), 10, 64)
	h.svc.AddStory(uint(stageID), c.PostForm("title"), c.PostForm("assignee"))
	redirectTo(c, fmt.Sprintf("/releases/%d", releaseID))
}

func (h *ReleaseHandler) DeleteStory(c *gin.Context) {
	storyID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	releaseID, _ := strconv.ParseUint(c.PostForm("release_id"), 10, 64)
	h.svc.DeleteStory(uint(storyID))
	redirectTo(c, fmt.Sprintf("/releases/%d", releaseID))
}

func (h *ReleaseHandler) UpdateStoryStatus(c *gin.Context) {
	storyID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	releaseID, _ := strconv.ParseUint(c.PostForm("release_id"), 10, 64)
	h.svc.UpdateStoryStatus(uint(storyID), c.PostForm("status"))
	redirectTo(c, fmt.Sprintf("/releases/%d", releaseID))
}

func (h *ReleaseHandler) UpdateStory(c *gin.Context) {
	storyID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	releaseID, _ := strconv.ParseUint(c.PostForm("release_id"), 10, 64)
	h.svc.UpdateStory(uint(storyID), c.PostForm("title"), c.PostForm("assignee"))
	redirectTo(c, fmt.Sprintf("/releases/%d", releaseID))
}

func (h *ReleaseHandler) CreateSlackUpdate(c *gin.Context) {
	stageID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	releaseID, _ := strconv.ParseUint(c.PostForm("release_id"), 10, 64)
	h.svc.AddSlackUpdate(
		uint(stageID),
		c.PostForm("channel"),
		c.PostForm("message"),
		c.PostForm("author"),
	)
	redirectTo(c, fmt.Sprintf("/releases/%d", releaseID))
}

func (h *ReleaseHandler) DeleteSlackUpdate(c *gin.Context) {
	slackID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	releaseID, _ := strconv.ParseUint(c.PostForm("release_id"), 10, 64)
	h.svc.DeleteSlackUpdate(uint(slackID))
	redirectTo(c, fmt.Sprintf("/releases/%d", releaseID))
}
