package service

import (
	"voidsounds/internal/domain"
	"voidsounds/internal/repository"
)

type EventService struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) GetAllEvents() (domain.Events, error) {
	return s.repo.GetAll()
}
