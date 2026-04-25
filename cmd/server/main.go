package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"voidsounds/internal/components"
	"voidsounds/internal/domain"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Роуты
	r.Get("/", homeHandler)
	r.Get("/events", eventsHandler)

	log.Println("VoidSounds запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	component := components.Home()
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	events := getMockEvents()
	component := components.Events(events)
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getMockEvents() domain.Events {
	return domain.Events{
		{
			ID:          1,
			Title:       "Шум и Выходки в баре «Подвал»",
			Description: "Сольный концерт группы Шум и Выходки. Инди-рок с мощным саундом и неожиданными каверами.",
			Date:        time.Date(2026, 5, 15, 20, 0, 0, 0, time.Local),
			Location:    "Бар «Подвал», Москва",
			Genre:       "Инди-рок",
			Price:       1500,
			Available:   87,
		},
		{
			ID:          2,
			Title:       "Электронная ночь на крыше",
			Description: "Три артиста электронной сцены. Живая электроника и визуальное шоу.",
			Date:        time.Date(2026, 5, 22, 22, 0, 0, 0, time.Local),
			Location:    "Крыша «Flora», Санкт-Петербург",
			Genre:       "Электронная",
			Price:       2000,
			Available:   45,
		},
	}
}
