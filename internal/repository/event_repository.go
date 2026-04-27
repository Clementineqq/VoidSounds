package repository

import (
	"time"

	"voidsounds/internal/domain"
)

type EventRepository interface {
	GetAll() (domain.Events, error)
	GetByID(id int) (*domain.Event, error)
	Create(event *domain.Event) error
	// GetByGenre, GetByCity добавим позже
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
		SELECT 
			id, title, description, date, city_id, address, 
			price, available, poster_url, organizer_id, 
			status, created_at, updated_at
		FROM events 
		WHERE status = 'published'
		ORDER BY date ASC`

	var events domain.Events
	err := DB.Select(&events, query)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *eventRepository) GetByID(id int) (*domain.Event, error) {
	// Пока оставим заглушкой
	return nil, nil
}

func (r *eventRepository) Create(event *domain.Event) error {
	// Реализуем позже при создании формы
	return nil
}

func getMockEvents() domain.Events {
	return domain.Events{
		{
			ID:          1,
			Title:       "Шум и Выходки в баре «Подвал»",
			Description: "Сольный концерт группы Шум и Выходки.",
			Date:        time.Date(2026, 5, 15, 20, 0, 0, 0, time.Local),
			Address:     "Бар «Подвал», Москва",
			Price:       1500,
			Available:   87,
			Status:      "published",
		},
	}
}
