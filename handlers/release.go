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
	releases, err := h.svc.All()
	if err != nil {
		c.String(500, "DB error: %s", err.Error())
		return
	}
	allProjects, activeProject := projectMeta(c)
	render(c, "releases", ReleasesData{
		Meta:     Meta{Title: "Releases", CurrentPage: "releases", AllProjects: allProjects, ActiveProject: activeProject},
		Releases: releases,
	})
}

func (h *ReleaseHandler) Create(c *gin.Context) {
	h.svc.Create(
		c.PostForm("name"),
		c.PostForm("description"),
		c.PostForm("status"),
		c.PostForm("target_date"),
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
	allProjects, activeProject := projectMeta(c)
	render(c, "release_detail", ReleaseDetailData{
		Meta:    Meta{Title: release.Name, CurrentPage: "releases", AllProjects: allProjects, ActiveProject: activeProject},
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
