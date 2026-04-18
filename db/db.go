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
		&models.SprintTask{},
		&models.SprintTaskComment{},
		&models.StandupEntry{},
		&models.Deadline{},
		&models.Meeting{},
		&models.ActionItem{},
		&models.DevTask{},
		&models.DevTaskComment{},
		&models.Release{},
		&models.ReleaseStage{},
		&models.ReleaseStory{},
		&models.ReleaseSlackUpdate{},
		&models.Project{},
		&models.TeamMember{},
		&models.SlackThread{},
	); err != nil {
		return err
	}
	return nil
}
