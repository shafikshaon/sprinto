package service

import (
	"strings"

	"sprinto/models"
	"sprinto/repository"
)

type SlackThreadService interface {
	All(tag string) ([]models.SlackThread, error)
	AllTags() ([]string, error)
	Create(messageLink, topic, summary, tags string, authorID *uint) error
	Update(id uint, messageLink, topic, summary, tags string, authorID *uint) error
	Delete(id uint) error
}

type slackThreadService struct{ repo repository.SlackThreadRepository }

func NewSlackThreadService(r repository.SlackThreadRepository) SlackThreadService {
	return &slackThreadService{repo: r}
}

func (s *slackThreadService) All(tag string) ([]models.SlackThread, error) {
	threads, err := s.repo.All(tag)
	if err != nil {
		return nil, err
	}
	for i := range threads {
		threads[i].Tags = splitTags(threads[i].TagCSV)
	}
	return threads, nil
}

func (s *slackThreadService) AllTags() ([]string, error) { return s.repo.AllTags() }

func (s *slackThreadService) Create(messageLink, topic, summary, tags string, authorID *uint) error {
	topic = strings.TrimSpace(topic)
	if topic == "" {
		return nil
	}
	return s.repo.Create(models.SlackThread{
		MessageLink: strings.TrimSpace(messageLink),
		Topic:       topic,
		Summary:     strings.TrimSpace(summary),
		TagCSV:      normaliseTags(tags),
		AuthorID:    authorID,
	})
}

func (s *slackThreadService) Update(id uint, messageLink, topic, summary, tags string, authorID *uint) error {
	return s.repo.Update(id,
		strings.TrimSpace(messageLink),
		strings.TrimSpace(topic),
		strings.TrimSpace(summary),
		normaliseTags(tags),
		authorID,
	)
}

func (s *slackThreadService) Delete(id uint) error { return s.repo.Delete(id) }

func normaliseTags(raw string) string {
	seen := map[string]bool{}
	var out []string
	for _, t := range strings.Split(raw, ",") {
		t = strings.ToLower(strings.TrimSpace(t))
		if t != "" && !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return strings.Join(out, ",")
}

func splitTags(csv string) []string {
	var tags []string
	for _, t := range strings.Split(csv, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}
