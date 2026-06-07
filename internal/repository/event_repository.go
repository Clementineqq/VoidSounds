package repository

import (
	"fmt"
	"time"
	"voidsounds/internal/domain"
)

type EventRepository interface {
	GetAll() (domain.Events, error)
	GetByID(id int) (*domain.Event, error)
	GetByOrganizerID(organizerID int) (domain.Events, error)
	Create(event *domain.Event) error
	Update(event *domain.Event) error
	Delete(id int) error
	BuyTicket(eventID, userID int) error
	GetTicketsByUserID(userID int) ([]domain.Ticket, error)
	GetAllCities() ([]domain.City, error)
	GetAllGenres() ([]domain.Genre, error)
	GetAllWithFilters(citySlug, genreSlug, search string) (domain.Events, error)
	GetGenresByEventID(eventID int) ([]domain.Genre, error)
	AssignGenresToEvent(eventID int, genreIDs []int) error
	GetAllEventsForAdmin() (domain.Events, error) // ← ДОБАВИТЬ
}

type eventRepository struct{}

func NewEventRepository() EventRepository {
	return &eventRepository{}
}

func (r *eventRepository) GetAll() (domain.Events, error) {
	if DB == nil {
		return getMockEvents(), nil
	}
	query := `SELECT id, title, description, date, city_id, address, price, available, poster_url, organizer_id, status, created_at, updated_at FROM events WHERE status = 'published' ORDER BY date ASC`
	var events domain.Events
	err := DB.Select(&events, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения событий: %w", err)
	}
	return events, nil
}

func (r *eventRepository) GetByID(id int) (*domain.Event, error) {
	if DB == nil {
		for _, event := range getMockEvents() {
			if event.ID == id {
				return &event, nil
			}
		}
		return nil, fmt.Errorf("событие с ID %d не найдено", id)
	}
	query := `SELECT e.id, e.title, e.description, e.date, e.city_id, e.address, 
							e.price, e.available, e.poster_url, e.organizer_id, e.status, 
							e.created_at, e.updated_at,
							u.name as organizer_name
					FROM events e 
					LEFT JOIN users u ON e.organizer_id = u.id 
					WHERE e.id = $1`
	var event domain.Event
	err := DB.Get(&event, query, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения события %d: %w", id, err)
	}
	return &event, nil
}

func (r *eventRepository) Create(event *domain.Event) error {
	if DB == nil {
		return fmt.Errorf("база данных не подключена")
	}
	query := `INSERT INTO events (title, description, date, city_id, address, price, available, poster_url, organizer_id, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, created_at, updated_at`
	err := DB.QueryRowx(query, event.Title, event.Description, event.Date, event.CityID, event.Address, event.Price, event.Available, event.PosterURL, event.OrganizerID, event.Status).Scan(&event.ID, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		return fmt.Errorf("ошибка создания события: %w", err)
	}
	return nil
}

func (r *eventRepository) Update(event *domain.Event) error {
	if DB == nil {
		return fmt.Errorf("база данных не подключена")
	}
	query := `UPDATE events SET title = $1, description = $2, date = $3, city_id = $4, address = $5, price = $6, available = $7, poster_url = $8, status = $9, updated_at = NOW() WHERE id = $10`
	result, err := DB.Exec(query, event.Title, event.Description, event.Date, event.CityID, event.Address, event.Price, event.Available, event.PosterURL, event.Status, event.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления события: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("событие с ID %d не найдено", event.ID)
	}
	return nil
}

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

func (r *eventRepository) GetByOrganizerID(organizerID int) (domain.Events, error) {
	if DB == nil {
		var events domain.Events
		for _, e := range getMockEvents() {
			if e.OrganizerID == organizerID {
				events = append(events, e)
			}
		}
		return events, nil
	}
	query := `SELECT id, title, description, date, city_id, address, price, available, poster_url, organizer_id, status, created_at, updated_at FROM events WHERE organizer_id = $1 ORDER BY date ASC`
	var events domain.Events
	err := DB.Select(&events, query, organizerID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения событий организатора: %w", err)
	}
	return events, nil
}

func (r *eventRepository) BuyTicket(eventID, userID int) error {
	if DB == nil {
		return fmt.Errorf("база данных не подключена")
	}
	tx, err := DB.Beginx()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback()
	var available int
	err = tx.Get(&available, "SELECT available FROM events WHERE id = $1 FOR UPDATE", eventID)
	if err != nil {
		return fmt.Errorf("ошибка проверки билетов: %w", err)
	}
	if available <= 0 {
		return fmt.Errorf("билетов больше нет")
	}
	_, err = tx.Exec("UPDATE events SET available = available - 1 WHERE id = $1", eventID)
	if err != nil {
		return fmt.Errorf("ошибка обновления события: %w", err)
	}
	_, err = tx.Exec(`INSERT INTO tickets (event_id, user_id, quantity, total_price, status) VALUES ($1, $2, 1, (SELECT price FROM events WHERE id = $1), 'paid')`, eventID, userID)
	if err != nil {
		return fmt.Errorf("ошибка создания билета: %w", err)
	}
	return tx.Commit()
}

func (r *eventRepository) GetTicketsByUserID(userID int) ([]domain.Ticket, error) {
	if DB == nil {
		return []domain.Ticket{}, nil
	}
	query := `SELECT t.id, t.event_id, t.user_id, t.quantity, t.total_price, t.purchase_date, t.status, e.title, e.date, e.address, e.poster_url, e.status as event_status FROM tickets t JOIN events e ON t.event_id = e.id WHERE t.user_id = $1 ORDER BY t.purchase_date DESC`
	var tickets []domain.Ticket
	err := DB.Select(&tickets, query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения билетов: %w", err)
	}
	return tickets, nil
}

func (r *eventRepository) GetAllWithFilters(citySlug, genreSlug, search string) (domain.Events, error) {
	if DB == nil {
		return getMockEvents(), nil
	}

	query := `SELECT DISTINCT e.id, e.title, e.description, e.date, e.city_id, e.address, 
                     e.price, e.available, e.poster_url, e.organizer_id, e.status, 
                     e.created_at, e.updated_at,
                     u.name as organizer_name
              FROM events e 
              LEFT JOIN users u ON e.organizer_id = u.id 
              LEFT JOIN event_genres eg ON e.id = eg.event_id
              LEFT JOIN genres g ON eg.genre_id = g.id
              WHERE e.status = 'published'`

	args := []interface{}{}
	paramIndex := 1

	if citySlug != "" {
		query += fmt.Sprintf(` AND (e.city_id = (SELECT id FROM cities WHERE slug = $%d) OR e.address ILIKE '%%' || (SELECT name FROM cities WHERE slug = $%d) || '%%')`, paramIndex, paramIndex)
		args = append(args, citySlug)
		paramIndex++
	}
	// Фильтр по жанру
	if genreSlug != "" {
		query += fmt.Sprintf(` AND g.slug = $%d`, paramIndex)
		args = append(args, genreSlug)
		paramIndex++
	}

	// Поиск
	if search != "" {
		query += fmt.Sprintf(` AND (e.title ILIKE $%d OR e.description ILIKE $%d)`, paramIndex, paramIndex+1)
		args = append(args, "%"+search+"%", "%"+search+"%")
		paramIndex += 2
	}

	query += " ORDER BY e.date ASC"

	var events domain.Events
	err := DB.Select(&events, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка фильтрации событий: %w", err)
	}

	return events, nil
}
func (r *eventRepository) GetAllCities() ([]domain.City, error) {
	if DB == nil {
		return []domain.City{}, nil
	}
	query := `SELECT id, name, slug FROM cities ORDER BY name ASC`
	var cities []domain.City
	err := DB.Select(&cities, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения городов: %w", err)
	}
	return cities, nil
}

func (r *eventRepository) GetAllGenres() ([]domain.Genre, error) {
	if DB == nil {
		return []domain.Genre{}, nil
	}
	query := `SELECT id, name, slug FROM genres ORDER BY name ASC`
	var genres []domain.Genre
	err := DB.Select(&genres, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения жанров: %w", err)
	}
	return genres, nil
}

func (r *eventRepository) GetGenresByEventID(eventID int) ([]domain.Genre, error) {
	if DB == nil {
		return []domain.Genre{}, nil
	}
	query := `SELECT g.id, g.name, g.slug FROM genres g JOIN event_genres eg ON g.id = eg.genre_id WHERE eg.event_id = $1 ORDER BY g.name`
	var genres []domain.Genre
	err := DB.Select(&genres, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения жанров: %w", err)
	}
	return genres, nil
}

func (r *eventRepository) AssignGenresToEvent(eventID int, genreIDs []int) error {
	if DB == nil {
		return fmt.Errorf("база данных не подключена")
	}
	// Сначала удаляем старые связи
	_, err := DB.Exec(`DELETE FROM event_genres WHERE event_id = $1`, eventID)
	if err != nil {
		return fmt.Errorf("ошибка удаления старых жанров: %w", err)
	}
	// Добавляем новые
	for _, genreID := range genreIDs {
		_, err := DB.Exec(`INSERT INTO event_genres (event_id, genre_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, eventID, genreID)
		if err != nil {
			return fmt.Errorf("ошибка добавления жанра %d: %w", genreID, err)
		}
	}
	return nil
}

func getMockEvents() domain.Events {
	return domain.Events{
		{ID: 1, Title: "Шум и Выходки в баре «Подвал»", Description: "Сольный концерт группы Шум и Выходки. Nintendo-core, чиптюн, эксперименты.", Date: time.Date(2026, 5, 15, 20, 0, 0, 0, time.Local), Address: "Самара, Бар Подвал", Price: 800, Available: 87, Status: "published", OrganizerID: 1},
		{ID: 2, Title: "Mitski", Description: "Лютый Арт перфоманс Митски в нашем доме!", Date: time.Date(2026, 5, 9, 22, 0, 0, 0, time.Local), Address: "Самара, Бар Подвал", Price: 9999, Available: 0, Status: "published", OrganizerID: 1},
	}
}

func (r *eventRepository) GetAllEventsForAdmin() (domain.Events, error) {
	if DB == nil {
		return getMockEvents(), nil
	}
	query := `SELECT e.id, e.title, e.description, e.date, e.city_id, e.address, 
                     e.price, e.available, e.poster_url, e.organizer_id, e.status, 
                     e.created_at, e.updated_at,
                     u.name as organizer_name
              FROM events e 
              LEFT JOIN users u ON e.organizer_id = u.id 
              ORDER BY e.date DESC`
	var events domain.Events
	err := DB.Select(&events, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения всех событий: %w", err)
	}
	return events, nil
}
