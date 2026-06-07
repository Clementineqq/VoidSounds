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
	return &EventService{repo: repo}
}

func (s *EventService) GetAllEvents() (domain.Events, error) {
	events, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("сервис: не удалось получить мероприятия: %w", err)
	}
	return events, nil
}

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

func (s *EventService) GetEventByIDForEdit(id int) (*domain.Event, error) {
	if id <= 0 {
		return nil, fmt.Errorf("неверный ID мероприятия: %d", id)
	}
	return s.repo.GetByID(id)
}

func (s *EventService) GetEventWithGenres(id int) (*domain.Event, error) {
	event, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	event.Genres, _ = s.repo.GetGenresByEventID(id)
	return event, nil
}

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

func (s *EventService) CreateEventWithGenres(event *domain.Event, genreIDs []int) error {
	if err := s.CreateEvent(event); err != nil {
		return err
	}
	if event.ID > 0 && len(genreIDs) > 0 {
		return s.repo.AssignGenresToEvent(event.ID, genreIDs)
	}
	return nil
}

func (s *EventService) UpdateEvent(eventID, organizerID int, req *domain.Event) error {
	existing, err := s.repo.GetByID(eventID)
	if err != nil || existing.OrganizerID != organizerID {
		return fmt.Errorf("мероприятие не найдено или недоступно")
	}
	existing.Title = req.Title
	existing.Description = req.Description
	existing.Date = req.Date
	existing.Address = req.Address
	existing.Price = req.Price
	existing.Available = req.Available
	existing.PosterURL = req.PosterURL
	existing.Status = req.Status
	return s.repo.Update(existing)
}

func (s *EventService) UpdateEventWithGenres(eventID, organizerID int, req *domain.Event, genreIDs []int) error {
	if err := s.UpdateEvent(eventID, organizerID, req); err != nil {
		return err
	}
	if len(genreIDs) > 0 {
		return s.repo.AssignGenresToEvent(eventID, genreIDs)
	}
	return nil
}

func (s *EventService) DeleteEvent(eventID, organizerID int) error {
	existing, err := s.repo.GetByID(eventID)
	if err != nil || existing.OrganizerID != organizerID {
		return fmt.Errorf("мероприятие не найдено или недоступно")
	}
	return s.repo.Delete(eventID)
}

func (s *EventService) BuyTicket(eventID, userID int) error {
	if eventID <= 0 || userID <= 0 {
		return fmt.Errorf("неверные параметры покупки")
	}
	event, err := s.repo.GetByID(eventID)
	if err != nil {
		return fmt.Errorf("мероприятие не найдено")
	}
	if event.OrganizerID == userID {
		return fmt.Errorf("организаторы не могут покупать билеты на свои мероприятия")
	}
	if event.Status != "published" {
		return fmt.Errorf("мероприятие не доступно для покупки")
	}
	if event.Available <= 0 {
		return fmt.Errorf("билеты закончились")
	}
	return s.repo.BuyTicket(eventID, userID)
}

func (s *EventService) GetOrganizerEvents(organizerID int) (domain.Events, error) {
	if organizerID <= 0 {
		return nil, fmt.Errorf("неверный ID организатора")
	}
	return s.repo.GetByOrganizerID(organizerID)
}

func (s *EventService) UpdateStatus(eventID, organizerID int, status string) error {
	existing, err := s.repo.GetByID(eventID)
	if err != nil || existing.OrganizerID != organizerID {
		return fmt.Errorf("мероприятие не найдено или нет прав")
	}
	existing.Status = status
	return s.repo.Update(existing)
}

func (s *EventService) GetUserTickets(userID int) ([]domain.Ticket, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("неверный ID пользователя")
	}
	return s.repo.GetTicketsByUserID(userID)
}

func (s *EventService) GetEventsWithFilters(citySlug, genreSlug, search string) (domain.Events, error) {
	return s.repo.GetAllWithFilters(citySlug, genreSlug, search)
}

func (s *EventService) GetAllCities() ([]domain.City, error) {
	return s.repo.GetAllCities()
}

func (s *EventService) GetAllGenres() ([]domain.Genre, error) {
	return s.repo.GetAllGenres()
}

func (s *EventService) GetGenresByEventID(eventID int) ([]domain.Genre, error) {
	return s.repo.GetGenresByEventID(eventID)
}

func (s *EventService) GetEventsByOrganizer(organizerID int) (domain.Events, error) {
	return s.repo.GetByOrganizerID(organizerID)
}

func (s *EventService) GetAllEventsForAdmin() (domain.Events, error) {
	return s.repo.GetAllEventsForAdmin()
}

func (s *EventService) DeleteEventAdmin(eventID int) error {
	return s.repo.Delete(eventID)
}

func (s *EventService) UpdateEventAdmin(eventID int, req *domain.Event) error {
	existing, err := s.repo.GetByID(eventID)
	if err != nil {
		return fmt.Errorf("мероприятие не найдено")
	}

	existing.Title = req.Title
	existing.Description = req.Description
	existing.Date = req.Date
	existing.Address = req.Address
	existing.Price = req.Price
	existing.Available = req.Available
	existing.PosterURL = req.PosterURL
	existing.Status = req.Status

	return s.repo.Update(existing)
}
