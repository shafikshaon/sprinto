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
	// ── Repositories ──────────────────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	sprintRepo := repository.NewSprintRepository(db)
	standupRepo := repository.NewStandupRepository(db)
	deadlineRepo := repository.NewDeadlineRepository(db)
	meetingRepo := repository.NewMeetingRepository(db)
	devTaskRepo := repository.NewDevTaskRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	teamMemberRepo := repository.NewTeamMemberRepository(db)
	slackThreadRepo := repository.NewSlackThreadRepository(db)
	stickyNoteRepo := repository.NewStickyNoteRepository(db)

	// ── Services ──────────────────────────────────────────────────
	authSvc := service.NewAuthService(userRepo)
	sprintSvc := service.NewSprintService(sprintRepo)
	standupSvc := service.NewStandupService(standupRepo)
	deadlineSvc := service.NewDeadlineService(deadlineRepo)
	meetingSvc := service.NewMeetingService(meetingRepo)
	devTaskSvc := service.NewDevTaskService(devTaskRepo)
	projectSvc := service.NewProjectService(projectRepo)
	teamMemberSvc := service.NewTeamMemberService(teamMemberRepo)
	slackThreadSvc := service.NewSlackThreadService(slackThreadRepo)
	stickyNoteSvc := service.NewStickyNoteService(stickyNoteRepo)

	// ── Handlers ──────────────────────────────────────────────────
	authH := handlers.NewAuthHandler(authSvc, teamMemberSvc)
	dashH := handlers.NewDashboardHandler(sprintSvc, standupSvc, deadlineSvc, devTaskSvc)
	sprintH := handlers.NewSprintHandler(sprintSvc)
	standupH := handlers.NewStandupHandler(standupSvc)
	deadlineH := handlers.NewDeadlineHandler(deadlineSvc)
	meetingH := handlers.NewMeetingHandler(meetingSvc)
	devTaskH := handlers.NewDevTaskHandler(devTaskSvc)
	projectH := handlers.NewProjectHandler(projectSvc, teamMemberSvc)
	teamH := handlers.NewTeamHandler(teamMemberSvc)
	slackH := handlers.NewSlackHandler(slackThreadSvc)
	notesH := handlers.NewStickyNoteHandler(stickyNoteSvc)

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

	// Sprint board + release (merged)
	protected.GET("/sprints", sprintH.List)
	protected.POST("/sprints", sprintH.CreateSprint)
	protected.POST("/sprints/:id/update", sprintH.UpdateSprint)
	protected.POST("/sprints/:id/delete", sprintH.DeleteSprint)
	protected.POST("/sprints/:id/activate", sprintH.ActivateSprint)
	protected.POST("/sprints/tasks", sprintH.CreateTask)
	protected.GET("/sprints/tasks/:id", sprintH.TaskDetail)
	protected.POST("/sprints/tasks/:id/update", sprintH.UpdateTask)
	protected.POST("/sprints/tasks/:id/delete", sprintH.DeleteTask)
	protected.POST("/sprints/tasks/:id/comments", sprintH.AddComment)
	protected.POST("/sprints/comments/:id/delete", sprintH.DeleteComment)
	protected.POST("/sprints/progress", sprintH.UpdateProgress)
	protected.POST("/sprints/release/update", sprintH.UpdateRelease)
	protected.POST("/sprints/stages", sprintH.CreateStage)
	protected.POST("/sprints/stages/:id/delete", sprintH.DeleteStage)
	protected.POST("/sprints/stages/:id/status", sprintH.UpdateStageStatus)
	protected.POST("/sprints/stages/:id/stories", sprintH.CreateStory)
	protected.POST("/sprints/stories/:id/delete", sprintH.DeleteStory)
	protected.POST("/sprints/stories/:id/status", sprintH.UpdateStoryStatus)
	protected.POST("/sprints/stories/:id/update", sprintH.UpdateStory)
	protected.POST("/sprints/stages/:id/slack", sprintH.CreateSlackUpdate)
	protected.POST("/sprints/slack/:id/delete", sprintH.DeleteSlackUpdate)

	protected.GET("/standups", standupH.List)
	protected.GET("/standups/pdf", standupH.PDF)
	protected.POST("/standups", standupH.Create)
	protected.POST("/standups/:id/update", standupH.Update)
	protected.POST("/standups/:id/delete", standupH.Delete)

	protected.GET("/deadlines", deadlineH.List)
	protected.POST("/deadlines", deadlineH.Create)
	protected.POST("/deadlines/:id/update", deadlineH.Update)
	protected.POST("/deadlines/:id/delete", deadlineH.Delete)

	protected.GET("/meetings", meetingH.List)
	protected.POST("/meetings", meetingH.Create)
	protected.POST("/meetings/:id/update", meetingH.Update)
	protected.POST("/meetings/:id/delete", meetingH.Delete)

	protected.GET("/devtasks", devTaskH.List)
	protected.POST("/devtasks", devTaskH.Create)
	protected.GET("/devtasks/:id", devTaskH.Detail)
	protected.POST("/devtasks/:id/update", devTaskH.Update)
	protected.POST("/devtasks/:id/delete", devTaskH.Delete)
	protected.POST("/devtasks/:id/comments", devTaskH.AddComment)
	protected.POST("/devtasks/comments/:id/delete", devTaskH.DeleteComment)

	protected.GET("/projects", projectH.List)
	protected.POST("/projects", projectH.Create)
	protected.POST("/projects/:id/update", projectH.Update)
	protected.POST("/projects/:id/delete", projectH.Delete)
	protected.POST("/projects/:id/members", projectH.AddMember)
	protected.POST("/projects/members/:id/remove", projectH.RemoveMember)
	protected.POST("/switch-project", projectH.SwitchProject)

	protected.GET("/team", teamH.List)
	protected.POST("/team", teamH.Create)
	protected.POST("/team/:id/update", teamH.Update)
	protected.POST("/team/:id/delete", teamH.Delete)

	protected.GET("/notes", notesH.List)
	protected.GET("/notes/new", notesH.New)
	protected.GET("/notes/:id/edit", notesH.EditPage)
	protected.POST("/notes", notesH.Create)
	protected.POST("/notes/:id/update", notesH.Update)
	protected.POST("/notes/:id/pin", notesH.TogglePin)
	protected.POST("/notes/:id/delete", notesH.Delete)

	protected.GET("/slack", slackH.List)
	protected.POST("/slack", slackH.Create)
	protected.POST("/slack/:id/update", slackH.Update)
	protected.POST("/slack/:id/delete", slackH.Delete)

	log.Printf("Sprinto running → http://localhost:%s", cfg.Port)
	r.Run(":" + cfg.Port)
}
