package repository

import (
	"sort"
	"strings"

	"gorm.io/gorm"
	"sprinto/models"
)

type SlackThreadRepository interface {
	All(tag string) ([]models.SlackThread, error)
	AllTags() ([]string, error)
	Create(t models.SlackThread) error
	Update(id uint, messageLink, topic, summary, tags string, authorID *uint) error
	Delete(id uint) error
}

type slackThreadRepo struct{ db *gorm.DB }

func NewSlackThreadRepository(db *gorm.DB) SlackThreadRepository {
	return &slackThreadRepo{db: db}
}

func (r *slackThreadRepo) All(tag string) ([]models.SlackThread, error) {
	var threads []models.SlackThread
	q := r.db.Preload("Author").Order("created_at DESC")
	if tag != "" {
		q = q.Where("tags LIKE ?", "%"+tag+"%")
	}
	return threads, q.Find(&threads).Error
}

func (r *slackThreadRepo) AllTags() ([]string, error) {
	var threads []models.SlackThread
	if err := r.db.Select("tags").Find(&threads).Error; err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var tags []string
	for _, t := range threads {
		for _, tag := range strings.Split(t.TagCSV, ",") {
			tag = strings.TrimSpace(tag)
			if tag != "" && !seen[tag] {
				seen[tag] = true
				tags = append(tags, tag)
			}
		}
	}
	sort.Strings(tags)
	return tags, nil
}

func (r *slackThreadRepo) Create(t models.SlackThread) error { return r.db.Create(&t).Error }

func (r *slackThreadRepo) Update(id uint, messageLink, topic, summary, tags string, authorID *uint) error {
	return r.db.Model(&models.SlackThread{}).Where("id = ?", id).Updates(map[string]interface{}{
		"message_link": messageLink,
		"topic":        topic,
		"summary":      summary,
		"tags":         tags,
		"author_id":    authorID,
	}).Error
}

func (r *slackThreadRepo) Delete(id uint) error {
	return r.db.Delete(&models.SlackThread{}, id).Error
}
