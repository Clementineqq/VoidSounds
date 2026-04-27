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
