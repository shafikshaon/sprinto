package handlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"sprinto/models"
	"sprinto/service"
)

type StandupsData struct {
	Meta        Meta
	Date        string
	DateRaw     string
	Standups    []models.StandupEntry
	RecentDates []service.DateNav
}

type StandupHandler struct {
	svc service.StandupService
}

func NewStandupHandler(svc service.StandupService) *StandupHandler {
	return &StandupHandler{svc: svc}
}

func (h *StandupHandler) List(c *gin.Context) {
	dateRaw := c.Query("date")
	if dateRaw == "" {
		dateRaw = time.Now().Format("2006-01-02")
	}
	t, _ := time.Parse("2006-01-02", dateRaw)
	entries, _ := h.svc.ByDate(dateRaw)
	recentDates, _ := h.svc.RecentDates(5)

	render(c, "standups", StandupsData{
		Meta:        Meta{Title: "Daily Standups", CurrentPage: "standups"},
		Date:        t.Format("Monday, January 2, 2006"),
		DateRaw:     dateRaw,
		Standups:    entries,
		RecentDates: recentDates,
	})
}

func (h *StandupHandler) Create(c *gin.Context) {
	h.svc.Add(
		c.PostForm("member"),
		c.PostForm("role"),
		c.PostForm("yesterday"),
		c.PostForm("today"),
		c.PostForm("blockers"),
		c.PostForm("status"),
	)
	redirectTo(c, "/standups")
}

func (h *StandupHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.svc.Remove(uint(id))
	redirectTo(c, "/standups")
}
