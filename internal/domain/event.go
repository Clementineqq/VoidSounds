package domain

import "time"

type Event struct {
	ID          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Date        time.Time `db:"date" json:"date"`
	Location    string    `db:"location" json:"location"`
	Genre       string    `db:"genre" json:"genre"`
	Price       int       `db:"price" json:"price"`
	Available   int       `db:"available" json:"available"`
	OrganizerID *int      `db:"organizer_id" json:"organizer_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type Events []Event
