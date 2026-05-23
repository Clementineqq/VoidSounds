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

	if event.Status != "published" {
		return nil, fmt.Errorf("мероприятие %d не опубликовано", id)
	}

	return event, nil
}

// CreateEvent - создаем мероприятие (будет использоваться организаторами)
func (s *EventService) CreateEvent(event *domain.Event) error {
	if event.Title == "" {
		return fmt.Errorf("название мероприятия не может быть пустым")
	}
	if event.Price < 0 {
		return fmt.Errorf("цена не может быть отрицательной")
	}
	if event.Available < 0 {
		return fmt.Errorf("количество билетов не может быть отрицательным")
	}

	if event.Status == "" {
		event.Status = "draft"
	}

	return s.repo.Create(event)
}

// BuyTicket - бизнес-логика покупки билета ← НОВЫЙ МЕТОД
func (s *EventService) BuyTicket(eventID, userID int) error {
	if eventID <= 0 || userID <= 0 {
		return fmt.Errorf("неверные параметры покупки")
	}

	event, err := s.repo.GetByID(eventID)
	if err != nil {
		return fmt.Errorf("мероприятие не найдено")
	}
	if event.Status != "published" {
		return fmt.Errorf("мероприятие не доступно для покупки")
	}
	if event.Available <= 0 {
		return fmt.Errorf("билеты закончились")
	}

	return s.repo.BuyTicket(eventID, userID)
}
