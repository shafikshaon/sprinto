package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"sprinto/config"
	appdb "sprinto/db"
	"sprinto/handlers"
	"sprinto/repository"
	"sprinto/service"
)

func main() {
	cfg := config.Load()

	// ── Database ──────────────────────────────────────────────────
	db := appdb.Connect(cfg.DatabaseURL)
	if err := appdb.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	if appdb.IsEmpty(db) {
		if err := appdb.Seed(db); err != nil {
			log.Fatalf("seed: %v", err)
		}
		log.Println("Database seeded with sample data")
	}

	// ── Repositories ──────────────────────────────────────────────
	sprintRepo := repository.NewSprintRepository(db)
	standupRepo := repository.NewStandupRepository(db)
	deadlineRepo := repository.NewDeadlineRepository(db)
	meetingRepo := repository.NewMeetingRepository(db)
	devTaskRepo := repository.NewDevTaskRepository(db)
	releaseRepo := repository.NewReleaseRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	teamMemberRepo := repository.NewTeamMemberRepository(db)

	// ── Services ──────────────────────────────────────────────────
	sprintSvc := service.NewSprintService(sprintRepo)
	standupSvc := service.NewStandupService(standupRepo)
	deadlineSvc := service.NewDeadlineService(deadlineRepo)
	meetingSvc := service.NewMeetingService(meetingRepo)
	devTaskSvc := service.NewDevTaskService(devTaskRepo)
	releaseSvc := service.NewReleaseService(releaseRepo)
	projectSvc := service.NewProjectService(projectRepo)
	teamMemberSvc := service.NewTeamMemberService(teamMemberRepo)

	// ── Handlers ──────────────────────────────────────────────────
	dashH := handlers.NewDashboardHandler(sprintSvc, standupSvc, deadlineSvc, devTaskSvc)
	sprintH := handlers.NewSprintHandler(sprintSvc)
	standupH := handlers.NewStandupHandler(standupSvc)
	deadlineH := handlers.NewDeadlineHandler(deadlineSvc)
	meetingH := handlers.NewMeetingHandler(meetingSvc)
	devTaskH := handlers.NewDevTaskHandler(devTaskSvc)
	releaseH := handlers.NewReleaseHandler(releaseSvc)
	projectH := handlers.NewProjectHandler(projectSvc, teamMemberSvc)
	teamH := handlers.NewTeamHandler(teamMemberSvc)

	// ── Router ────────────────────────────────────────────────────
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(handlers.ProjectMiddleware(projectSvc))

	r.GET("/", dashH.Get)

	r.GET("/sprints", sprintH.List)
	r.POST("/sprints/tasks", sprintH.CreateTask)
	r.GET("/sprints/tasks/:id", sprintH.TaskDetail)
	r.POST("/sprints/tasks/:id/delete", sprintH.DeleteTask)
	r.POST("/sprints/tasks/:id/comments", sprintH.AddComment)
	r.POST("/sprints/comments/:id/delete", sprintH.DeleteComment)
	r.POST("/sprints/progress", sprintH.UpdateProgress)

	r.GET("/standups", standupH.List)
	r.POST("/standups", standupH.Create)
	r.POST("/standups/:id/delete", standupH.Delete)

	r.GET("/deadlines", deadlineH.List)
	r.POST("/deadlines", deadlineH.Create)
	r.POST("/deadlines/:id/delete", deadlineH.Delete)

	r.GET("/meetings", meetingH.List)
	r.POST("/meetings", meetingH.Create)
	r.POST("/meetings/:id/delete", meetingH.Delete)

	r.GET("/devtasks", devTaskH.List)
	r.POST("/devtasks", devTaskH.Create)
	r.GET("/devtasks/:id", devTaskH.Detail)
	r.POST("/devtasks/:id/delete", devTaskH.Delete)
	r.POST("/devtasks/:id/comments", devTaskH.AddComment)
	r.POST("/devtasks/comments/:id/delete", devTaskH.DeleteComment)

	r.GET("/releases", releaseH.List)
	r.POST("/releases", releaseH.Create)
	r.POST("/releases/:id/delete", releaseH.Delete)
	r.GET("/releases/:id", releaseH.Detail)
	r.POST("/releases/:id/stages", releaseH.CreateStage)
	r.POST("/releases/stages/:id/delete", releaseH.DeleteStage)
	r.POST("/releases/stages/:id/status", releaseH.UpdateStageStatus)
	r.POST("/releases/stages/:id/stories", releaseH.CreateStory)
	r.POST("/releases/stories/:id/delete", releaseH.DeleteStory)
	r.POST("/releases/stories/:id/status", releaseH.UpdateStoryStatus)
	r.POST("/releases/stages/:id/slack", releaseH.CreateSlackUpdate)
	r.POST("/releases/slack/:id/delete", releaseH.DeleteSlackUpdate)

	r.GET("/projects", projectH.List)
	r.POST("/projects", projectH.Create)
	r.POST("/projects/:id/delete", projectH.Delete)
	r.POST("/projects/:id/members", projectH.AddMember)
	r.POST("/projects/members/:id/remove", projectH.RemoveMember)
	r.POST("/switch-project", projectH.SwitchProject)

	r.GET("/team", teamH.List)
	r.POST("/team", teamH.Create)
	r.POST("/team/:id/delete", teamH.Delete)

	log.Printf("Sprinto running → http://localhost:%s", cfg.Port)
	r.Run(":" + cfg.Port)
}
