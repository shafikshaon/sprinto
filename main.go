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

	// ── Services ──────────────────────────────────────────────────
	sprintSvc := service.NewSprintService(sprintRepo)
	standupSvc := service.NewStandupService(standupRepo)
	deadlineSvc := service.NewDeadlineService(deadlineRepo)
	meetingSvc := service.NewMeetingService(meetingRepo)
	devTaskSvc := service.NewDevTaskService(devTaskRepo)

	// ── Handlers ──────────────────────────────────────────────────
	dashH := handlers.NewDashboardHandler(sprintSvc, standupSvc, deadlineSvc, devTaskSvc)
	sprintH := handlers.NewSprintHandler(sprintSvc)
	standupH := handlers.NewStandupHandler(standupSvc)
	deadlineH := handlers.NewDeadlineHandler(deadlineSvc)
	meetingH := handlers.NewMeetingHandler(meetingSvc)
	devTaskH := handlers.NewDevTaskHandler(devTaskSvc)

	// ── Router ────────────────────────────────────────────────────
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/", dashH.Get)

	r.GET("/sprints", sprintH.List)
	r.POST("/sprints/tasks", sprintH.CreateTask)
	r.POST("/sprints/tasks/:id/delete", sprintH.DeleteTask)
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
	r.POST("/devtasks/:id/delete", devTaskH.Delete)

	log.Printf("Sprinto running → http://localhost:%s", cfg.Port)
	r.Run(":" + cfg.Port)
}
