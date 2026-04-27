package repository

import (
	"time"

	"voidsounds/internal/domain"
)

type EventRepository interface {
	GetAll() (domain.Events, error)
	GetByID(id int) (*domain.Event, error)
	// Create, Update, Delete добавим позже
}

type eventRepository struct{}

func NewEventRepository() EventRepository {
	return &eventRepository{}
}

func (r *eventRepository) GetAll() (domain.Events, error) {
	if DB == nil {
		return getMockEvents(), nil
	}

	query := `
		SELECT id, title, description, date, location, genre, price, available, 
		       organizer_id, created_at 
		FROM events 
		ORDER BY date ASC`

	var events domain.Events
	err := DB.Select(&events, query)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *eventRepository) GetByID(id int) (*domain.Event, error) {
	return nil, nil // пока не реализуем
}

func getMockEvents() domain.Events {
	return domain.Events{
		{
			ID:          1,
			Title:       "Шум и Выходки в баре «Подвал»",
			Description: "Сольный концерт группы Шум и Выходки.",
			Date:        time.Date(2026, 5, 15, 20, 0, 0, 0, time.Local),
			Location:    "Бар «Подвал», Москва",
			Genre:       "Инди-рок",
			Price:       1500,
			Available:   87,
		},
	}
}
