package domain

import "time"

type Event struct {
	ID          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Date        time.Time `db:"date" json:"date"`
	CityID      *int      `db:"city_id" json:"city_id"`
	Address     string    `db:"address" json:"address"`
	Price       int       `db:"price" json:"price"`
	Available   int       `db:"available" json:"available"`
	PosterURL   *string   `db:"poster_url" json:"poster_url"` // ← изменили на *string
	OrganizerID int       `db:"organizer_id" json:"organizer_id"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type Events []Event

// Вспомогательная структура для жанров
type Genre struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Slug string `db:"slug" json:"slug"`
}

// Ticket - купленный билет (расширенный для отображения)
type Ticket struct {
	ID           int       `db:"id" json:"id"`
	EventID      int       `db:"event_id" json:"event_id"`
	UserID       int       `db:"user_id" json:"user_id"`
	Quantity     int       `db:"quantity" json:"quantity"`
	TotalPrice   int       `db:"total_price" json:"total_price"`
	PurchaseDate time.Time `db:"purchase_date" json:"purchase_date"`
	Status       string    `db:"status" json:"status"`

	// Данные мероприятия (для отображения)
	EventTitle     string    `db:"title" json:"event_title"`
	EventDate      time.Time `db:"date" json:"event_date"`
	EventAddress   string    `db:"address" json:"event_address"`
	EventPosterURL *string   `db:"poster_url" json:"event_poster_url"`
	EventStatus    string    `db:"event_status" json:"event_status"` // ← ДОБАВИТЬ

}
