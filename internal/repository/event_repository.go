package repository

import (
	"fmt"
	"time"
	"voidsounds/internal/domain"
)

type EventRepository interface {
	GetAll() (domain.Events, error)
	GetByID(id int) (*domain.Event, error)
	Create(event *domain.Event) error
	Update(event *domain.Event) error
	Delete(id int) error
}

type eventRepository struct{}

func NewEventRepository() EventRepository {
	return &eventRepository{}
}

// GetAll - получаем все опубликованные мероприятия
func (r *eventRepository) GetAll() (domain.Events, error) {
	if DB == nil {
		// Если БД не подключена, возвращаем мок-данные для тестирования
		return getMockEvents(), nil
	}

	query := `
        SELECT 
            id, title, description, date, city_id, address,
            price, available, poster_url, organizer_id,
            status, created_at, updated_at
        FROM events
        WHERE status = 'published'
        ORDER BY date ASC
    `

	var events domain.Events
	err := DB.Select(&events, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения событий: %w", err)
	}

	return events, nil
}

// GetByID - получаем мероприятие по ID (даже если не опубликовано, для организатора)
func (r *eventRepository) GetByID(id int) (*domain.Event, error) {
	if DB == nil {
		// Для разработки: ищем в мок-данных
		for _, event := range getMockEvents() {
			if event.ID == id {
				return &event, nil
			}
		}
		return nil, fmt.Errorf("событие с ID %d не найдено", id)
	}

	query := `
        SELECT 
            id, title, description, date, city_id, address,
            price, available, poster_url, organizer_id,
            status, created_at, updated_at
        FROM events
        WHERE id = $1
    `

	var event domain.Event
	err := DB.Get(&event, query, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения события %d: %w", id, err)
	}

	return &event, nil
}

// Create - создаем новое мероприятие
func (r *eventRepository) Create(event *domain.Event) error {
	if DB == nil {
		return fmt.Errorf("база данных не подключена")
	}

	query := `
        INSERT INTO events (
            title, description, date, city_id, address,
            price, available, poster_url, organizer_id, status
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, created_at, updated_at
    `

	err := DB.QueryRowx(
		query,
		event.Title,
		event.Description,
		event.Date,
		event.CityID,
		event.Address,
		event.Price,
		event.Available,
		event.PosterURL,
		event.OrganizerID,
		event.Status,
	).Scan(&event.ID, &event.CreatedAt, &event.UpdatedAt)

	if err != nil {
		return fmt.Errorf("ошибка создания события: %w", err)
	}

	return nil
}

// Update - обновляем существующее мероприятие
func (r *eventRepository) Update(event *domain.Event) error {
	if DB == nil {
		return fmt.Errorf("база данных не подключена")
	}

	query := `
        UPDATE events SET
            title = $1,
            description = $2,
            date = $3,
            city_id = $4,
            address = $5,
            price = $6,
            available = $7,
            poster_url = $8,
            status = $9,
            updated_at = NOW()
        WHERE id = $10
    `

	result, err := DB.Exec(
		query,
		event.Title,
		event.Description,
		event.Date,
		event.CityID,
		event.Address,
		event.Price,
		event.Available,
		event.PosterURL,
		event.Status,
		event.ID,
	)

	if err != nil {
		return fmt.Errorf("ошибка обновления события: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("событие с ID %d не найдено", event.ID)
	}

	return nil
}

// Delete - удаляем мероприятие
func (r *eventRepository) Delete(id int) error {
	if DB == nil {
		return fmt.Errorf("база данных не подключена")
	}

	query := `DELETE FROM events WHERE id = $1`
	result, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления события: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("событие с ID %d не найдено", id)
	}

	return nil
}

// getMockEvents - временные тестовые данные
func getMockEvents() domain.Events {
	return domain.Events{
		{
			ID:          1,
			Title:       "Шум и Выходки в баре «Подвал»",
			Description: "Сольный концерт группы Шум и Выходки. Nintendo-core, чиптюн, эксперименты.",
			Date:        time.Date(2026, 5, 15, 20, 0, 0, 0, time.Local),
			Address:     "Бар «Подвал», Мытищи",
			Price:       800,
			Available:   87,
			Status:      "published",
			OrganizerID: 1,
		},
		{
			ID:          2,
			Title:       "Mitski",
			Description: "Лютый Арт перфоманс Митски в нашем доме!",
			Date:        time.Date(2026, 5, 9, 22, 0, 0, 0, time.Local),
			Address:     "квартира жинки любимой, где-то в Самаре",
			Price:       9999,
			Available:   0,
			Status:      "published",
			OrganizerID: 1,
		},
	}
}
