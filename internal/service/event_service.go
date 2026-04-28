package service

import (
	"fmt"
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

// GetAllEvents - получаем все мероприятия
func (s *EventService) GetAllEvents() (domain.Events, error) {
	events, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("сервис: не удалось получить мероприятия: %w", err)
	}
	return events, nil
}

// GetEventByID - получаем мероприятие по ID с проверками
func (s *EventService) GetEventByID(id int) (*domain.Event, error) {
	if id <= 0 {
		return nil, fmt.Errorf("неверный ID мероприятия: %d", id)
	}

	event, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("сервис: мероприятие %d не найдено: %w", id, err)
	}

	// Проверяем, что мероприятие опубликовано (или показываем только для организатора)
	if event.Status != "published" {
		return nil, fmt.Errorf("мероприятие %d не опубликовано", id)
	}

	return event, nil
}

// CreateEvent - создаем мероприятие (будет использоваться организаторами)
func (s *EventService) CreateEvent(event *domain.Event) error {
	// Валидация данных
	if event.Title == "" {
		return fmt.Errorf("название мероприятия не может быть пустым")
	}
	if event.Price < 0 {
		return fmt.Errorf("цена не может быть отрицательной")
	}
	if event.Available < 0 {
		return fmt.Errorf("количество билетов не может быть отрицательным")
	}

	// Если статус не указан, ставим "draft"
	if event.Status == "" {
		event.Status = "draft"
	}

	return s.repo.Create(event)
}
