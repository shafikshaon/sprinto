package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"sprinto/models"
)

func Connect(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("db.Connect: %v\nMake sure PostgreSQL is running.\nCreate the DB with: createdb sprinto", err)
	}
	return db
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Sprint{},
		&models.Task{},
		&models.TaskComment{},
		&models.ReleaseStage{},
		&models.ReleaseSlackUpdate{},
		&models.StandupEntry{},
		&models.Deadline{},
		&models.Meeting{},
		&models.ActionItem{},
		&models.Project{},
		&models.TeamMember{},
		&models.SlackThread{},
		&models.StickyNote{},
	); err != nil {
		return err
	}
	// Drop legacy NOT NULL + FK on release_stages.release_id (merged into sprints).
	db.Exec(`ALTER TABLE release_stages DROP CONSTRAINT IF EXISTS fk_releases_stages`)
	db.Exec(`ALTER TABLE release_stages ALTER COLUMN release_id DROP NOT NULL`)
	return nil
}
