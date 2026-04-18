package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sprinto/models"
	"sprinto/service"
)

type TeamData struct {
	Meta    Meta
	Members []models.TeamMember
}

type TeamHandler struct {
	svc service.TeamMemberService
}

func NewTeamHandler(svc service.TeamMemberService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

func (h *TeamHandler) List(c *gin.Context) {
	members, err := h.svc.All()
	if err != nil {
		c.String(500, "DB error: %s", err.Error())
		return
	}
	allProjects, activeProject := projectMeta(c)
	render(c, "team", TeamData{
		Meta:    Meta{Title: "Team", CurrentPage: "team", ActionLabel: "Add Member", AllProjects: allProjects, ActiveProject: activeProject},
		Members: members,
	})
}

func (h *TeamHandler) Create(c *gin.Context) {
	h.svc.Create(c.PostForm("name"), c.PostForm("role"), c.PostForm("email"))
	redirectTo(c, "/team")
}

func (h *TeamHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Delete(uint(id))
	redirectTo(c, "/team")
}
