package domain

import "time"

type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
	Genre       string    `json:"genre"`
	Price       int       `json:"price"`
	Available   int       `json:"available"`
	OrganizerID int       `json:"organizer_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Events []Event
