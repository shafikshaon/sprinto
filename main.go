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
	userRepo := repository.NewUserRepository(db)
	sprintRepo := repository.NewSprintRepository(db)
	standupRepo := repository.NewStandupRepository(db)
	deadlineRepo := repository.NewDeadlineRepository(db)
	meetingRepo := repository.NewMeetingRepository(db)
	devTaskRepo := repository.NewDevTaskRepository(db)
	releaseRepo := repository.NewReleaseRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	teamMemberRepo := repository.NewTeamMemberRepository(db)
	slackThreadRepo := repository.NewSlackThreadRepository(db)

	// ── Services ──────────────────────────────────────────────────
	authSvc := service.NewAuthService(userRepo)
	sprintSvc := service.NewSprintService(sprintRepo)
	standupSvc := service.NewStandupService(standupRepo)
	deadlineSvc := service.NewDeadlineService(deadlineRepo)
	meetingSvc := service.NewMeetingService(meetingRepo)
	devTaskSvc := service.NewDevTaskService(devTaskRepo)
	releaseSvc := service.NewReleaseService(releaseRepo)
	projectSvc := service.NewProjectService(projectRepo)
	teamMemberSvc := service.NewTeamMemberService(teamMemberRepo)
	slackThreadSvc := service.NewSlackThreadService(slackThreadRepo)

	// ── Handlers ──────────────────────────────────────────────────
	authH := handlers.NewAuthHandler(authSvc)
	dashH := handlers.NewDashboardHandler(sprintSvc, standupSvc, deadlineSvc, devTaskSvc)
	sprintH := handlers.NewSprintHandler(sprintSvc)
	standupH := handlers.NewStandupHandler(standupSvc)
	deadlineH := handlers.NewDeadlineHandler(deadlineSvc)
	meetingH := handlers.NewMeetingHandler(meetingSvc)
	devTaskH := handlers.NewDevTaskHandler(devTaskSvc)
	releaseH := handlers.NewReleaseHandler(releaseSvc)
	projectH := handlers.NewProjectHandler(projectSvc, teamMemberSvc)
	teamH := handlers.NewTeamHandler(teamMemberSvc)
	slackH := handlers.NewSlackHandler(slackThreadSvc)

	// ── Router ────────────────────────────────────────────────────
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Public routes (no auth)
	r.GET("/login", authH.LoginPage)
	r.POST("/login", authH.Login)
	r.GET("/register", authH.RegisterPage)
	r.POST("/register", authH.Register)
	r.POST("/logout", authH.Logout)

	// Protected routes
	protected := r.Group("/")
	protected.Use(handlers.LoadUserMiddleware(authSvc))
	protected.Use(handlers.AuthRequiredMiddleware())
	protected.Use(handlers.ProjectMiddleware(projectSvc))

	protected.GET("/", dashH.Get)

	protected.GET("/sprints", sprintH.List)
	protected.POST("/sprints/tasks", sprintH.CreateTask)
	protected.GET("/sprints/tasks/:id", sprintH.TaskDetail)
	protected.POST("/sprints/tasks/:id/update", sprintH.UpdateTask)
	protected.POST("/sprints/tasks/:id/delete", sprintH.DeleteTask)
	protected.POST("/sprints/tasks/:id/comments", sprintH.AddComment)
	protected.POST("/sprints/comments/:id/delete", sprintH.DeleteComment)
	protected.POST("/sprints/progress", sprintH.UpdateProgress)

	protected.GET("/standups", standupH.List)
	protected.POST("/standups", standupH.Create)
	protected.POST("/standups/:id/update", standupH.Update)
	protected.POST("/standups/:id/delete", standupH.Delete)

	protected.GET("/deadlines", deadlineH.List)
	protected.POST("/deadlines", deadlineH.Create)
	protected.POST("/deadlines/:id/update", deadlineH.Update)
	protected.POST("/deadlines/:id/delete", deadlineH.Delete)

	protected.GET("/meetings", meetingH.List)
	protected.POST("/meetings", meetingH.Create)
	protected.POST("/meetings/:id/delete", meetingH.Delete)

	protected.GET("/devtasks", devTaskH.List)
	protected.POST("/devtasks", devTaskH.Create)
	protected.GET("/devtasks/:id", devTaskH.Detail)
	protected.POST("/devtasks/:id/update", devTaskH.Update)
	protected.POST("/devtasks/:id/delete", devTaskH.Delete)
	protected.POST("/devtasks/:id/comments", devTaskH.AddComment)
	protected.POST("/devtasks/comments/:id/delete", devTaskH.DeleteComment)

	protected.GET("/releases", releaseH.List)
	protected.POST("/releases", releaseH.Create)
	protected.POST("/releases/:id/delete", releaseH.Delete)
	protected.GET("/releases/:id", releaseH.Detail)
	protected.POST("/releases/:id/stages", releaseH.CreateStage)
	protected.POST("/releases/stages/:id/delete", releaseH.DeleteStage)
	protected.POST("/releases/stages/:id/status", releaseH.UpdateStageStatus)
	protected.POST("/releases/stages/:id/stories", releaseH.CreateStory)
	protected.POST("/releases/stories/:id/delete", releaseH.DeleteStory)
	protected.POST("/releases/stories/:id/status", releaseH.UpdateStoryStatus)
	protected.POST("/releases/stories/:id/update", releaseH.UpdateStory)
	protected.POST("/releases/stages/:id/slack", releaseH.CreateSlackUpdate)
	protected.POST("/releases/slack/:id/delete", releaseH.DeleteSlackUpdate)

	protected.GET("/projects", projectH.List)
	protected.POST("/projects", projectH.Create)
	protected.POST("/projects/:id/update", projectH.Update)
	protected.POST("/projects/:id/delete", projectH.Delete)
	protected.POST("/projects/:id/members", projectH.AddMember)
	protected.POST("/projects/members/:id/remove", projectH.RemoveMember)
	protected.POST("/switch-project", projectH.SwitchProject)

	protected.GET("/team", teamH.List)
	protected.POST("/team", teamH.Create)
	protected.POST("/team/:id/delete", teamH.Delete)

	protected.GET("/slack", slackH.List)
	protected.POST("/slack", slackH.Create)
	protected.POST("/slack/:id/update", slackH.Update)
	protected.POST("/slack/:id/delete", slackH.Delete)

	log.Printf("Sprinto running → http://localhost:%s", cfg.Port)
	r.Run(":" + cfg.Port)
}
