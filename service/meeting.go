package service

import (
	"strings"

	"sprinto/models"
	"sprinto/repository"
)

type MeetingService interface {
	All(projectID uint) ([]models.Meeting, error)
	Add(title, date, attendeesCSV, notes string, projectID uint) error
	Remove(id uint) error
}

type meetingService struct{ repo repository.MeetingRepository }

func NewMeetingService(r repository.MeetingRepository) MeetingService {
	return &meetingService{repo: r}
}

func (s *meetingService) All(projectID uint) ([]models.Meeting, error) {
	meetings, err := s.repo.All(projectID)
	if err != nil {
		return nil, err
	}
	for i := range meetings {
		if meetings[i].AttendeeCSV != "" {
			for _, a := range strings.Split(meetings[i].AttendeeCSV, ",") {
				if t := strings.TrimSpace(a); t != "" {
					meetings[i].Attendees = append(meetings[i].Attendees, t)
				}
			}
		}
	}
	return meetings, nil
}

func (s *meetingService) Add(title, date, attendeesCSV, notes string, projectID uint) error {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(date) == "" {
		return nil
	}
	var attendees []string
	for _, a := range strings.Split(attendeesCSV, ",") {
		if t := strings.TrimSpace(a); t != "" {
			attendees = append(attendees, t)
		}
	}
	return s.repo.Create(models.Meeting{
		ProjectID:   projectID,
		Title:       strings.TrimSpace(title),
		Date:        strings.TrimSpace(date),
		AttendeeCSV: strings.Join(attendees, ","),
		Notes:       strings.TrimSpace(notes),
	})
}

func (s *meetingService) Remove(id uint) error { return s.repo.Delete(id) }
