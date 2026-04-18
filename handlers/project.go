package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sprinto/models"
	"sprinto/service"
)

// ── Middleware ────────────────────────────────────────────────────────────────

// ProjectMiddleware fetches all projects + resolves the active project from
// the "active_project" cookie, then stores both in the Gin context so every
// handler can call projectMeta(c) without needing its own ProjectService.
func ProjectMiddleware(svc service.ProjectService) gin.HandlerFunc {
	return func(c *gin.Context) {
		projects, _ := svc.AllWithMembers()
		var active *models.Project
		if cookie, err := c.Cookie("active_project"); err == nil && cookie != "" {
			id, _ := strconv.ParseUint(cookie, 10, 64)
			for i := range projects {
				if projects[i].ID == uint(id) {
					active = &projects[i]
					break
				}
			}
		}
		c.Set("all_projects", projects)
		c.Set("active_project", active)
		c.Next()
	}
}

// ── Handler ───────────────────────────────────────────────────────────────────

type ProjectsData struct {
	Meta       Meta
	Projects   []models.Project
	AllMembers []models.TeamMember
}

type ProjectHandler struct {
	svc    service.ProjectService
	teamSvc service.TeamMemberService
}

func NewProjectHandler(svc service.ProjectService, teamSvc service.TeamMemberService) *ProjectHandler {
	return &ProjectHandler{svc: svc, teamSvc: teamSvc}
}

func (h *ProjectHandler) List(c *gin.Context) {
	projects, err := h.svc.AllWithMembers()
	if err != nil {
		c.String(500, "DB error: %s", err.Error())
		return
	}
	members, _ := h.teamSvc.All()
	allProjects, activeProject, currentUser := projectMeta(c)
	render(c, "projects", ProjectsData{
		Meta:       Meta{Title: "Projects", CurrentPage: "projects", ActionLabel: "New Project", AllProjects: allProjects, ActiveProject: activeProject, CurrentUser: currentUser},
		Projects:   projects,
		AllMembers: members,
	})
}

func (h *ProjectHandler) Create(c *gin.Context) {
	h.svc.Create(c.PostForm("name"), c.PostForm("description"))
	redirectTo(c, "/projects")
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Update(uint(id), c.PostForm("name"), c.PostForm("description"))
	redirectTo(c, "/projects")
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Delete(uint(id))
	redirectTo(c, "/projects")
}

func (h *ProjectHandler) AddMember(c *gin.Context) {
	projectID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	memberID, _ := strconv.ParseUint(c.PostForm("member_id"), 10, 64)
	h.svc.AddMember(uint(projectID), uint(memberID))
	redirectTo(c, "/projects")
}

func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	projectID, _ := strconv.ParseUint(c.PostForm("project_id"), 10, 64)
	memberID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.RemoveMember(uint(projectID), uint(memberID))
	redirectTo(c, "/projects")
}

func (h *ProjectHandler) SwitchProject(c *gin.Context) {
	id := c.PostForm("project_id")
	if id == "" {
		c.SetCookie("active_project", "", -1, "/", "", false, false)
	} else {
		c.SetCookie("active_project", id, 86400*30, "/", "", false, false)
	}
	ref := c.Request.Referer()
	if ref == "" {
		ref = "/"
	}
	redirectTo(c, ref)
}
