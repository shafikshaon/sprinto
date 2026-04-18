package handlers

import (
	"time"

	"github.com/gin-gonic/gin"

	"sprinto/models"
	"sprinto/service"
)

type DashboardData struct {
	Meta              Meta
	Sprint            models.Sprint
	TodayStandup      *models.StandupEntry
	Deadlines         []models.Deadline
	CriticalDeadlines int
	OpenDevTasks      int
}

type DashboardHandler struct {
	sprints   service.SprintService
	standups  service.StandupService
	deadlines service.DeadlineService
	devTasks  service.DevTaskService
}

func NewDashboardHandler(
	sp service.SprintService,
	su service.StandupService,
	dl service.DeadlineService,
	dt service.DevTaskService,
) *DashboardHandler {
	return &DashboardHandler{sprints: sp, standups: su, deadlines: dl, devTasks: dt}
}

func (h *DashboardHandler) Get(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)
	sprint, err := h.sprints.ActiveSprint(projectID)
	if err != nil {
		c.String(500, "DB error: %s", err.Error())
		return
	}
	todayStandups, _ := h.standups.ByDate(time.Now().Format("2006-01-02"), projectID)
	allDeadlines, _ := h.deadlines.All(projectID)
	allTasks, _ := h.devTasks.All(projectID)

	open, critical := 0, 0
	for _, t := range allTasks {
		if t.Status != "Done" {
			open++
		}
	}
	for _, d := range allDeadlines {
		if d.DaysLeft <= 3 {
			critical++
		}
	}
	cap := 3
	if len(allDeadlines) < cap {
		cap = len(allDeadlines)
	}
	var todayStandup *models.StandupEntry
	if len(todayStandups) > 0 {
		todayStandup = &todayStandups[0]
	}

	sprintLabel := sprint.Name
	if sprint.StartDate != "" && sprint.EndDate != "" {
		sprintLabel += " · " + sprint.StartDate + " – " + sprint.EndDate
	}
	allProjects, activeProject, currentUser := projectMeta(c)

	render(c, "dashboard", DashboardData{
		Meta:              Meta{Title: "Dashboard", CurrentPage: "dashboard", SprintLabel: sprintLabel, AllProjects: allProjects, ActiveProject: activeProject, CurrentUser: currentUser},
		Sprint:            sprint,
		TodayStandup:      todayStandup,
		Deadlines:         allDeadlines[:cap],
		CriticalDeadlines: critical,
		OpenDevTasks:      open,
	})
}
